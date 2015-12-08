package lib
import (
	"errors"
	"fmt"
)

// Goroutine票池的接口。
type GoTickets interface {
	// 取得一张票
	Take()

	//归还一张票
	Return()

	// 票池是否被激活
	Active() bool

	// 票的总数
	Total() uint32

	// 剩余的票数
	Remainder() uint32
}

// Goroutine票池的实现。
type myGoTickets struct {
	total    uint32    // 票的总数。
	ticketCh chan byte // 票的容器。
	active   bool      // 票池是否已被激活。
}

func NewGoTickets(total uint32) (GoTickets, error) {
	gt := myGoTickets{}
	if !gt.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

// 初始化
func (gt *myGoTickets) init(total uint32) bool {
	// 如果已经激活则返回false, 则不需要继续初始化
	if gt.active {
		return false
	}
	// 如果初始化的票的总数为零值,则不可以初始化
	if total == 0 {
		return false
	}

	// 创建一个含有total个的缓冲channel
	ch := make(chan byte, total)

	n := int(total)

	for i := 0; i < n; i++ {
		ch <- 1
	}

	gt.ticketCh = ch
	gt.total = total
	gt.active = true

	return true
}

func (gt *myGoTickets) Take() {
	<- gt.ticketCh
}

func (gt *myGoTickets) Return() {
	gt.ticketCh <- 1
}

func (gt *myGoTickets) Active() bool {
	return gt.active
}

func (gt *myGoTickets) Total() uint32 {
	return gt.total
}

func (gt *myGoTickets) Remainder() uint32 {
	return uint32(len(gt.ticketCh))
}


