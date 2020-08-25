package model

const IVA = 1.16
const COPUSD = 3800
type Product interface {
	GetPriceCOP() float64
	GetPriceUSD() float64
	GetName() string
	GetId() string
}

type BasicProduct struct {
	Name     string
	PriceCOP float64
	TypeBasic bool
	Id string
}
type NormalProduct struct {
	Name     string
	PriceCOP float64
	TypeNormal bool
	Id string
}

func (p BasicProduct) GetPriceCOP() float64 {
	return p.PriceCOP
}

func (p BasicProduct) GetPriceUSD() float64 {
	return p.PriceCOP / COPUSD
}
func (p BasicProduct) GetName() string  {
	return p.Name
}

func (p BasicProduct) GetId() string  {
	return p.Id
}

func (p NormalProduct) GetPriceCOP() float64  {
	return p.PriceCOP * IVA
}

func (p NormalProduct) GetPriceUSD() float64  {
	return (p.PriceCOP / COPUSD)* IVA
}

func (p NormalProduct) GetName() string  {
	return p.Name
}
func (p NormalProduct) GetId() string  {
	return p.Id
}



