package monkey

/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/

func networkCompute(a, b int) (int, error) {
	c := a + b
	return c, nil
}

func Compute(a, b int) (int, error) {
	sum, err := networkCompute(a, b)
	/*
		业务逻辑
	*/
	return sum, err
}

type Computer struct {
}

func (c *Computer) NetworkCompute(a, b int) (int, error) {
	sum := a + b
	return sum, nil
}

func (c *Computer) Compute(a, b int) (int, error) {
	sum, err := c.NetworkCompute(a, b)
	/*
		业务逻辑
	*/
	return sum, err
}
