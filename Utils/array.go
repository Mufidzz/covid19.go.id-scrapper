package Utils

import s "../Struct"

func ReverseNewsArray(a []s.News) []s.News {
	var nA []s.News

	for i := range a {
		nA = append(nA, a[len(a)-1-i])
	}

	return nA
}

func ReverseHoaxArray(a []s.Hoax) []s.Hoax {
	var nA []s.Hoax

	for i := range a {
		nA = append(nA, a[len(a)-1-i])
	}

	return nA
}

func ReverseEducationArray(a []s.Education) []s.Education {
	var nA []s.Education

	for i := range a {
		nA = append(nA, a[len(a)-1-i])
	}

	return nA
}

func ReverseProtocolArray(a []s.Protocol) []s.Protocol {
	var nA []s.Protocol

	for i := range a {
		nA = append(nA, a[len(a)-1-i])
	}

	return nA
}
