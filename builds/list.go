/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats"
	"github.com/spf13/cobra"
)

type Build struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func List(cmd *cobra.Command, args []string) {
	var builds []Build

	nuri, _ := cmd.Flags().GetString("nats")
	env, _ := cmd.Flags().GetString("env")

	nc, err := nats.Connect(nuri)
	if err != nil {
		panic(err)
	}

	q := []byte(`{"name": "` + env + `"}`)
	if env == "" {
		q = nil
	}

	msg, err := nc.Request("service.find", q, time.Second)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(msg.Data, &builds)
	if err != nil {
		panic(err)
	}

	for _, b := range builds {
		fmt.Printf("%s  %s  %s\n", b.Name, b.Status, b.ID)
	}
}
