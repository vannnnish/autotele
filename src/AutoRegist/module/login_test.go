/**
 * Created by vannnnish on 2018/6/5.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package module

import (
	"time"
	"testing"
	"fmt"
)

func TestTicker(t *testing.T) {
	for {
		select {
		case <-time.Tick(5 * time.Second):
			fmt.Println("是不是五秒执行一次")
		case <-time.Tick(1 * time.Second):
			fmt.Println("一秒一次")
		}
	}
}
