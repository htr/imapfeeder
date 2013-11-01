package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"code.google.com/p/go-imap/go1/imap"
)

type ImapMsg struct {
	folder      string
	extraCopies []string
	body        io.Reader
}

type ImapSession struct {
	client  *imap.Client
	isGmail bool
}

/* TODO: chan'ize the whole thing

func imapSink(msgs <-chan ImapMsg) {

	existingMBs := map[string]bool{}

	for msg := range msgs {
		if _, ok := existingMBs[msg.folder]; !ok {
			err := ctx.Imap.CreateFolder(msg.folder)
			if err != nil {
				fmt.Println("unable to create folder", msg.folder, " => ", err)
				continue
			} else {
				existingMBs[msg.folder] = true
			}
		}

		ctx.Imap.Append(msg.folder, msg.extraCopies, msg.body)
	}
}

*/

func imapConnect(addr, user, pass string) (*ImapSession, error) {
	c, err := imap.DialTLS(addr, nil)

	if err != nil {
		return nil, err
	}

	_, err = c.Login(user, pass)

	if err != nil {
		return nil, err
	}

	imapSession := &ImapSession{client: c}

	if c.Caps["X-GM-EXT-1"] {
		imapSession.isGmail = true
	}

	return imapSession, nil
}

func (i *ImapSession) IsConnected() bool {
	return i.client != nil && i.client.State() != imap.Closed
}

func (i *ImapSession) Close() {
	i.client.Logout(1 * time.Second)
}

func (i *ImapSession) imapCleanup(folders []string) error {

	c := i.client

	for _, mbox := range folders {
		cmd, err := imap.Wait(c.Select(mbox, false))
		if err != nil {
			fmt.Println("unable to select", mbox, "=>", err, i.client.State())
			continue
		}
		fmt.Println("cleaning up", mbox)
		yesterday := time.Now().Add(-1 * 24 * time.Hour)

		cmd = imapMust(imap.Wait(c.UIDSearch("SEEN BEFORE " + yesterday.Format("02-Jan-2006") + " NOT FLAGGED")))

		toDelete, _ := imap.NewSeqSet("")
		toDelete.AddNum(cmd.Data[0].SearchResults()...)
		if !toDelete.Empty() {
			fmt.Println("deleting...", toDelete)
			if i.isGmail {
				imapMust(imap.Wait(c.UIDStore(toDelete, "X-GM-LABELS", imap.NewFlagSet(`\Trash`))))
			} else {
				imapMust(imap.Wait(c.UIDStore(toDelete, "+FLAGS.SILENT", imap.NewFlagSet(`\Deleted`))))
			}
			imapMust(imap.Wait(c.Expunge(nil)))
		}
	}

	return nil
}

func imapMust(cmd *imap.Command, err error) *imap.Command {
	if err != nil {
		panic(err)
	}
	return cmd
}

func (i *ImapSession) CreateFolder(folder string) error {
	c := i.client

	mbox := folder

	if _, err := imap.Wait(c.Create(mbox)); err != nil {
		if rsp, ok := err.(imap.ResponseError); ok && rsp.Status == imap.NO {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (i *ImapSession) Append(folder string, extraCopies []string, msg io.Reader) error {
	c := i.client

	mbox := folder

	buf := new(bytes.Buffer)
	buf.ReadFrom(msg)

	now := time.Now()
	cmd, err := imap.Wait(c.Append(mbox, nil, &now, imap.NewLiteral(buf.Bytes())))
	if err != nil {
		fmt.Println(err)
		return err
	}

	rsp, err := cmd.Result(imap.OK)
	if err != nil {
		fmt.Println(err)
		return err
	}

	uid := imap.AsNumber(rsp.Fields[len(rsp.Fields)-1])

	set, _ := imap.NewSeqSet("")
	set.AddNum(uid)

	imapMust(imap.Wait(c.Select(mbox, false)))

	for _, mb := range extraCopies {
		if i.isGmail {
			imapMust(imap.Wait(c.UIDStore(set, "X-GM-LABELS", imap.NewFlagSet(mb))))
		} else {
			imapMust(c.UIDCopy(set, mb))
		}
	}

	return nil
}
