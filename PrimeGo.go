package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func main() {

	fmt.Println("PrimeGo by Larin Alexsandr v. 2018.8.25")
	fmt.Println("alexsandr@larin.name")
	fmt.Println("https://computerraru.ru")
	fmt.Println()

	if len(os.Args) != 3 {
		PrintError("Необходимо указать команду и экспоненту\r\nНапример, primego.exe number 615558040052742364849\r\nНапример, primego.exe mersenne 1277")
		return
	}

	n := os.Args[2]
	if !IsInteger(n) {
		PrintError("Второй аргумент должен быть числом")
		return
	}

	if !Pretest() {
		PrintError("Тестирование не пройдено")
		return
	}

	command := strings.ToLower(os.Args[1])
	switch command {
	case "number":
		FactorNumber(n)
		return
	case "mersenne":
		FactorMersenne(n)
		return
	}

	PrintError("Необходимо указать команду number или mersenne")

}

// JSONSaveStruct is ...
type JSONSaveStruct struct {
	Divider string
}

// Pretest is ...
func Pretest() bool {
	one := big.NewInt(1)
	two := big.NewInt(2)
	p, _ := new(big.Int).SetString("1277", 10)
	pow := new(big.Int).Exp(two, p, nil)
	m := new(big.Int).Sub(pow, one)
	n1277, _ := new(big.Int).SetString("2601983048666099770481310081841021384653815561816676201329778087600902014918340074503059860433081046210605403488570251947845891562080866227034976651419330190731032377347305086443295837415395887618239855136922452802923419286887119716740625346109565072933087221327790207134604146257063901166556207972729700461767055550785130256674608872183239507219512717434046725178680177638925792182271", 10)
	return m.Cmp(n1277) == 0
}

// LoadSaveFromFile is ...
func LoadSaveFromFile(path string) *big.Int {
	zero := big.NewInt(0)
	file, err := os.Open(path)
	if err != nil {
		return zero
	}
	defer file.Close()
	content := make([]byte, 64)
	for {
		bytes, err := file.Read(content)
		if err == io.EOF {
			break
		}
		jsn := string(content[:bytes])
		var jsonStruct JSONSaveStruct
		json.Unmarshal([]byte(jsn), &jsonStruct)
		divider, err1 := new(big.Int).SetString(jsonStruct.Divider, 10)
		if !err1 {
			return zero
		}
		return divider
	}
	return zero
}

// SavePointToFile is ...
func SavePointToFile(p *big.Int, k *big.Int) {
	file, err := os.Create(p.String() + ".json")
	if err != nil {
		fmt.Println("Не удалось сохранить в файл", err)
		fmt.Println(err)
	} else {
		defer file.Close()
		file.WriteString("{\r\n")
		file.WriteString("\t\"divider\": \"" + k.String() + "\"\r\n")
		file.WriteString("}\r\n")
	}
}

// SaveDividerToFile is ...
func SaveDividerToFile(p *big.Int, k *big.Int, d *big.Int, bit string) {
	file, err := os.Create(p.String() + "." + k.String() + ".json")
	if err != nil {
		fmt.Println("Не удалось сохранить в файл", err)
		fmt.Println(err)
	} else {
		defer file.Close()
		file.WriteString("{\r\n")
		file.WriteString("\t\"divider\": \"" + d.String() + "\",\r\n")
		file.WriteString("\t\"k\": \"" + k.String() + "\",\r\n")
		file.WriteString("\t\"bit\": \"" + bit + "\",\r\n")
		file.WriteString("}\r\n")
	}
}

// FloatToString is ...
func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 16, 64)
}

// BigIntToFloat64 is ...
func BigIntToFloat64(b *big.Int) float64 {
	f, _ := strconv.ParseFloat(b.String(), 64)
	return f
}

// BigIntLog2 is ...
func BigIntLog2(b *big.Int) float64 {
	f, _ := strconv.ParseFloat(b.String(), 64)
	return math.Log2(f)
}

// PrintError is ...
func PrintError(message string) {
	fmt.Println("Ошибка\r\n------\r\n" + message)
}

// IsInteger is ...
func IsInteger(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// FactorNumber is ...
func FactorNumber(n string) {
	start := time.Now()
	fmt.Println("Факторизуем число")
	fmt.Println()

	zero := big.NewInt(0)
	two := big.NewInt(2)
	three := big.NewInt(3)
	five := big.NewInt(5)
	six := big.NewInt(6)
	prime, _ := new(big.Int).SetString(n, 10)
	p05 := new(big.Int).Sqrt(prime)
	p05bits := BigIntLog2(p05)
	divider := LoadSaveFromFile(prime.String() + ".json")

	fmt.Println("Предварительные вычисления")
	fmt.Println("--------------------------")
	fmt.Println("Zero      = ", zero.String())
	fmt.Println("Two       = ", two.String())
	fmt.Println("Three     = ", three.String())
	fmt.Println("Five      = ", five.String())
	fmt.Println("Six       = ", six.String())
	fmt.Println("Prime     = ", prime.String())
	fmt.Println("Prime^0.5 = ", p05.String())
	fmt.Println("Bits      = ", p05bits)
	fmt.Println()

	if prime.Cmp(two) == -1 {
		fmt.Println("Мы проверяем на простоту только числа, которые больше 1")
		return
	}

	// Нужно проверить - если число равно 2, то оно просто, иначе проверить его на делимость на 2
	if prime.Cmp(two) == 0 {
		fmt.Println("Число официально простое!")
		return
	}
	reminder := new(big.Int).Mod(prime, two)
	if reminder.Cmp(zero) == 0 {
		SaveDividerToFile(prime, zero, two, "")
		return
	}

	// Нужно проверить - если число равно 3, то оно простое, иначе проверить его на делимость на 3
	if prime.Cmp(three) == 0 {
		fmt.Println("Число официально простое!")
		return
	}
	reminder = new(big.Int).Mod(prime, three)
	if reminder.Cmp(zero) == 0 {
		SaveDividerToFile(prime, zero, three, "")
		return
	}

	if divider.Cmp(three) > 0 {
		fmt.Println("Восстановление из файла")
		fmt.Println("-----------------------")
		fmt.Println("Успешно d = ", divider)
	}

	fmt.Println("Факторизация\r\n------------")

	reminder = new(big.Int).Mod(prime, two)
	if reminder.Cmp(zero) == 0 {
		fmt.Println("Делитель найден ", two.String())
		return
	}

	reminder = new(big.Int).Mod(prime, three)
	if reminder.Cmp(zero) == 0 {
		fmt.Println("Делитель найден ", three.String())
		return
	}

	c := 0

	for {
		c++
		divider = new(big.Int).Add(divider, two)
		if c%98765432 == 0 {
			if divider.Cmp(p05) == 1 {
				fmt.Println("Расчет окончен")
				return
			}
			elapsed := time.Since(start)
			bits := FloatToString(BigIntLog2(divider))
			fmt.Println("Претендент " + divider.String() + ", " + elapsed.String() + ", " + bits)
			start = time.Now()
			c = 0
			file, err := os.Create(prime.String() + ".json")
			if err != nil {
				fmt.Println("Не удалось сохранить в файл")
				fmt.Println(err)
			} else {
				defer file.Close()
				file.WriteString("{\r\n")
				file.WriteString("\t\"d\": \"" + divider.String() + "\"\r\n")
				file.WriteString("}\r\n")
			}

		}

		r := new(big.Int).Mod(prime, divider)
		if r.Cmp(zero) == 0 {
			SaveDividerToFile(prime, zero, divider, "")
		}

	}

}

// FactorMersenne is ...
func FactorMersenne(n string) {
	start := time.Now()
	fmt.Println("Факторизуем кандидата в Мерсенны")
	fmt.Println()
	cnt := 0
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	p, _ := new(big.Int).SetString(n, 10)
	p2 := new(big.Int).Mul(p, two)
	pow := new(big.Int).Exp(two, p, nil)
	m := new(big.Int).Sub(pow, one)
	sqrt := new(big.Int).Sqrt(m)
	kmax := new(big.Int).Div(new(big.Int).Sub(sqrt, one), p2)
	k := LoadSaveFromFile(p.String() + ".json")

	if k.Cmp(zero) > 0 {
		fmt.Println("Восстановление сохранения из файла")
		fmt.Println("----------------------------------")
		fmt.Println("Успешно")
		fmt.Println()
	}

	fmt.Println("Предварительные вычисления")
	fmt.Println("--------------------------")
	fmt.Println("Zero = " + zero.String())
	fmt.Println("One  = " + one.String())
	fmt.Println("Two  = " + two.String())
	fmt.Println("P    = " + p.String())
	fmt.Println("2*P  = " + p2.String())
	fmt.Println()
	fmt.Println("2^P:")
	fmt.Println(pow.String())
	fmt.Println()
	fmt.Println("2^P-1:")
	fmt.Println(m.String())
	fmt.Println()
	fmt.Println("SQRT(2^P-1):")
	fmt.Println(sqrt.String())
	fmt.Println()
	fmt.Println("Max bit SQRT(2^P-1) = " + FloatToString(BigIntLog2(sqrt)))
	fmt.Println("Kmin = " + k.String())
	fmt.Println("Kmax:")
	fmt.Println(kmax.String())
	fmt.Println()

	for {
		cnt++
		k = new(big.Int).Add(k, one)
		d := new(big.Int).Add(new(big.Int).Mul(p2, k), one)
		r := new(big.Int).Mod(m, d)
		if r.Cmp(zero) == 0 {
			if d.Cmp(m) == 0 {
				fmt.Println("Расчет окончен")
				return
			}
			bit := FloatToString(BigIntLog2(d))
			fmt.Println()
			fmt.Println("Найден делитель    : " + d.String())
			fmt.Println("K                  : " + k.String())
			fmt.Println("Бит                : " + bit)
			fmt.Println()
			SaveDividerToFile(p, k, d, bit)
		}
		if cnt%48654321 == 0 {
			if d.Cmp(sqrt) == 1 {
				fmt.Println("Расчет окончен")
				return
			}
			bit := FloatToString(BigIntLog2(d))
			elapsed := time.Since(start)
			fmt.Println("Претендент " + d.String() + ", бит " + bit + ", " + elapsed.String())
			start = time.Now()
			cnt = 0
			SavePointToFile(p, k)
		}
	}

}

/*

func FloatToString(input_num float64) string {
	//return strconv.FormatFloat(input_num, 'f', -1, 64)
	return strconv.FormatFloat(input_num, 'f', 16, 64)
}

func yglog2(x *big.Int) int {
	i := 0
	t := big.NewInt(2)
	for {
		x = new(big.Int).Div(x, t)
		i++
		if x.Cmp(big.NewInt(1)) == 0 {
			return i
		}
	}
}

func bigintlog2(b *big.Int) float64 {
	f, _ := strconv.ParseFloat(b.String(), 64)
	return math.Log2(f)
}

type SavedK struct {
	K string
}

func loadfromfile(f string) *big.Int {
	file, err := os.Open(f)
	if err != nil {
		return big.NewInt(0)
	}
	defer file.Close()
	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		j := string(data[:n])
		var savedk SavedK
		json.Unmarshal([]byte(j), &savedk)
		p, _ := new(big.Int).SetString(savedk.K, 10)
		return p
	}
	return big.NewInt(0)
}

*/
