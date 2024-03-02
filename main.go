package main

import (
	"fmt"
	"sort"
	"strconv"
	// печатание чека, давайте сохраним природу ?
)

type language string

var (
	languageRussian language = "russian"
	languageEnglish language = "english"
)

type Clients struct {
	amountOfWithdrawnMoney int
	balance                int
	CardNumber             int
	Name                   string
	PinCode                int
}

type ATM struct {
	clientID        int
	balance         int
	banknoteToCount map[int]int
	Clients         map[int]Clients
	language        language
}

func (a *ATM) ChangeLanguage(enterNumberOfLanguage int) {
	switch enterNumberOfLanguage {
	case 1:
		a.language = languageRussian
	case 2:
		a.language = languageEnglish
	}
}

func (a *ATM) CheckBalance() {
	fmt.Println("Ваш баланс:", a.Clients[a.clientID].balance)
}

func (a *ATM) Deposit(deposit int) {
	tempDeposit, estimateWithdrawanBanknoteAmount := deposit, 0
	tepmStorageBanknotes := []int{5000, 2000, 1000, 500, 100}

	for _, banknote := range tepmStorageBanknotes {
		estimateWithdrawanBanknoteAmount = tempDeposit / banknote
		a.banknoteToCount[banknote] += estimateWithdrawanBanknoteAmount
		tempDeposit -= estimateWithdrawanBanknoteAmount * banknote
	}

	if v, ok := a.Clients[a.clientID]; ok {
		v.balance += deposit
		a.Clients[a.clientID] = v
	}
	a.balance += deposit
	fmt.Println("Идёт распознавание купюр\nВы успешно пополнили:", deposit)
	fmt.Println(a.banknoteToCount)
}

func (a *ATM) DisplayAvalibleBanknotes() string {
	var displayBanknotes string

	for banknote := range a.banknoteToCount {
		displayBanknotes += strconv.Itoa(banknote) + " "
	}

	fmt.Println("Выберите сумму для снятия\n1 - ввести другую сумму\n", displayBanknotes) //выводим существующие банкноты
	return displayBanknotes
}

func (a *ATM) CheckPinForWithdraw(pinCode int) bool {
	fmt.Println("Введите пинкод")
	fmt.Scan(&pinCode)

	if pinCode != a.Clients[a.clientID].PinCode {
		fmt.Println("Не правильный пинкод")
		return false
	}
	return true
}

func (a *ATM) DisplayAvalibleExchange(withdrawal, сhoosingExchange int, displayBanknotes string) int {
	if withdrawal >= 500 && withdrawal <= 120000 {
		fmt.Println("Купюры какого размена желаете ?\n1 - Всё равно") //

		for banknote := range a.banknoteToCount {
			if withdrawal >= banknote && withdrawal < 6000 { // 6000, т.к. это максимальная сумма по минимальным купюром в данном кейсе
				displayBanknotes += strconv.Itoa(banknote) + " "
			} else if withdrawal <= banknote*60 && withdrawal >= 6000 {
				displayBanknotes += strconv.Itoa(banknote) + " "
			}
		}
		fmt.Println(displayBanknotes)
		fmt.Scan(&сhoosingExchange)
	} else {
		сhoosingExchange = 5000
	}
	return сhoosingExchange
}

func (a *ATM) CheckWithdrawAmount(withdrawal, ATMAmountForCheck int) bool {
	switch {
	case withdrawal%100 != 0:
		fmt.Println("Банкомата не выдаёт такие суммы, выберите другую сумму")
		return false
	case withdrawal > a.Clients[a.clientID].balance:
		fmt.Println("У вас недостаточно денег для снятия такой суммы")
		return false
	case withdrawal+a.Clients[a.clientID].amountOfWithdrawnMoney > 300000:
		fmt.Println("Максимальная сумма снятия с карты в день - 300 000")
		return false
	case withdrawal > ATMAmountForCheck:
		fmt.Println("Выберите сумму поменьше")
		return false
	}

	if _, ex := a.banknoteToCount[100]; !ex { // if нет 100, идёт проверка вводимого числа на кратность 100, но не 500 или вывод < 500
		if withdrawal%100 == 0 && withdrawal%500 != 0 || withdrawal < 500 {
			fmt.Println("Введите другую сумму")
			return false
		}
	}
	return true
}

func (a *ATM) WithdrawOneBanknote(chousenBanknote int) {
	totalAmountOneBanknoteMoney := chousenBanknote * a.banknoteToCount[chousenBanknote]
	if !a.CheckWithdrawAmount(chousenBanknote, totalAmountOneBanknoteMoney) {
		return
	}

	a.banknoteToCount[chousenBanknote] -= 1
	a.balance -= chousenBanknote
	if v, ex := a.Clients[a.clientID]; ex {
		v.balance -= chousenBanknote
		a.Clients[a.clientID] = v
	}
	fmt.Println("Вы успешно сняли:", chousenBanknote)
}

func (a *ATM) Withdraw() {
	var withdrawal, chousenBanknote, сhoosingExchange, ATMAmountForCheck, countBanknotes, pinCode int
	banknotesStorage := []int{5000, 2000, 1000, 500, 100}

	displayBanknotes := a.DisplayAvalibleBanknotes()
	fmt.Scan(&chousenBanknote)

	if !a.CheckPinForWithdraw(pinCode) {
		return
	}

	if chousenBanknote != 1 { //если выбор банкноты один из предложенных
		a.WithdrawOneBanknote(chousenBanknote)
		return
	}

	fmt.Println("Введите сумму для снятия наличных\nВведите сумму до 300000\nДоступны номиналы купюр:", displayBanknotes)
	fmt.Scan(&withdrawal)
	estimatedWithdrawal := withdrawal

	displayBanknotes = "" // переиспользуем переменную

	сhoosingExchange = a.DisplayAvalibleExchange(withdrawal, сhoosingExchange, displayBanknotes)

	for banknote, ammount := range a.banknoteToCount { // сумирование предполагаемого вывода включительно до выбранной банкноты
		if banknote <= сhoosingExchange {
			ATMAmountForCheck += banknote * ammount
		}
	}

	if !a.CheckWithdrawAmount(withdrawal, ATMAmountForCheck) {
		return
	}

	sort.Sort(sort.Reverse(sort.IntSlice(banknotesStorage)))

	estimatedWithdrawalBanknotes := map[int]int{
		100:  a.banknoteToCount[100],
		500:  a.banknoteToCount[500],
		1000: a.banknoteToCount[1000],
		2000: a.banknoteToCount[2000],
		5000: a.banknoteToCount[5000],
	}

	estimatedWithdrawalBanknotes, withdrawal = WithdrawMoneyFromATM(banknotesStorage, сhoosingExchange, withdrawal, estimatedWithdrawalBanknotes)

	if !a.CheckWithdrawBanknotes(estimatedWithdrawalBanknotes, сhoosingExchange, countBanknotes) {
		return
	}

	a.UpdateATMOptions(estimatedWithdrawalBanknotes, estimatedWithdrawal)
	a.UpdateClientOptions(estimatedWithdrawal)
	fmt.Println("Выдача купюр. . .\nВы успешно сняли", estimatedWithdrawal)
}

func (a *ATM) UpdateATMOptions(estimatedWithdrawalBanknotes map[int]int, estimatedWithdrawal int) {
	for banknote, v := range estimatedWithdrawalBanknotes { //окончательная передача информации о выводе денег
		a.banknoteToCount[banknote] = v
	}

	a.balance -= estimatedWithdrawal // обновляем баланс банкомата
}

func (a *ATM) CheckWithdrawBanknotes(estimatedWithdrawalBanknotes map[int]int, сhoosingExchange, countBanknotes int) bool {
	for banknote, v := range estimatedWithdrawalBanknotes { //  проверка количества снимаемых банкнот
		if banknote <= сhoosingExchange {
			countBanknotes += a.banknoteToCount[banknote] - v
			if countBanknotes > 60 {
				fmt.Println("банкомат не может выдать такое количество купюр")
				return false
			}
		}
	}
	return true
}

func (a *ATM) UpdateClientOptions(estimatedWithdrawal int) {
	if v, ok := a.Clients[a.clientID]; ok {
		v.amountOfWithdrawnMoney += estimatedWithdrawal
		v.balance -= estimatedWithdrawal
		a.Clients[a.clientID] = v
	}
}

func (a *ATM) checkBanknotesForExistence() bool {
	var countBeforeATMBlock int
	for _, v := range a.banknoteToCount {
		if v == 0 { // Проверка количества банкнот каждого номинала
			countBeforeATMBlock++
		}
	}
	if countBeforeATMBlock >= 2 {
		fmt.Println("Банкомат заблокирован")
		return false
	}
	return true
}

func WithdrawMoneyFromATM(banknotesStorage []int, withdrawal, сhoosingExchange int, estimatedWithdrawalAmount map[int]int) (map[int]int, int) {
	for _, banknote := range banknotesStorage { // процесс снятия денег с банкомата
		if withdrawal != 0 && banknote <= сhoosingExchange {

			tempEstimatedWithdrawalAmount := estimatedWithdrawalAmount[banknote]
			estimatedWithdrawalAmount[banknote] -= withdrawal / banknote

			if estimatedWithdrawalAmount[banknote] < 0 {
				estimatedWithdrawalAmount[banknote] = 0
			}
			if withdrawal/banknote != 0 {
				withdrawal -= (tempEstimatedWithdrawalAmount - estimatedWithdrawalAmount[banknote]) * banknote
			}
		}
	}
	return estimatedWithdrawalAmount, withdrawal
}

func (a *ATM) actionsOfATM() {
	fmt.Println("1 - снять наличные \n2 - пополнить карту \n3 - посмотреть баланс \n4 - поменять язык \n5 - вытащить карту")
	var enterNum int
	fmt.Scan(&enterNum)
	switch enterNum {

	case 1:
		a.Withdraw()
		if !a.checkBanknotesForExistence() {
			return
		}
		getBack()
		a.actionsOfATM()

	case 2:
		var deposit int
		fmt.Println("Введите сумму для пополнения")
		fmt.Scan(&deposit)
		if deposit%100 != 0 {
			fmt.Println("Банкомата не принимает такие суммы, выберите другую сумму")
			return
		}
		a.Deposit(deposit)
		getBack()
		a.actionsOfATM()

	case 3:
		a.CheckBalance()
		getBack()
		a.actionsOfATM()

	case 4:
		var enterNumberOfLanguage int
		fmt.Println("1 - поменять на русский\n2 - поменять на английский")
		fmt.Scan(&enterNumberOfLanguage)
		a.ChangeLanguage(enterNumberOfLanguage)
		getBack()
		a.actionsOfATM()

	case 5:
		fmt.Println("Заберите карту")
		a.theFirstEntering()
	}
}

func (a *ATM) checkTheFirstEnteringPincode() bool {
	var pin, countWrongPin int

	for a.Clients[a.clientID].PinCode != pin {
		if countWrongPin == 3 {
			fmt.Println("Ваша карта забанена")
			return false
		}

		fmt.Println("Введите пинкод")
		fmt.Scan(&pin)
		countWrongPin++
	}
	return true
}

func (a *ATM) theFirstEntering() {
	fmt.Println(a.balance)
	fmt.Println("Вставьте карту или приложите к NFC чипу \n1 - вставить карту или приложите к NFC чипу \n2 - поменять язык \n3 - выйти")
	var displayNum, cardNum int

	fmt.Scan(&displayNum)
	switch displayNum {

	case 1:
		fmt.Println("Введите номер карты")
		fmt.Scan(&cardNum)
		fmt.Println("Идёт считывание карты")

		if !a.checkClientForExistance(cardNum) {
			a.createNewClient(cardNum)
		} else {
			if !a.checkTheFirstEnteringPincode() {
				return
			}
		}

		fmt.Println("Добро пожаловать ʕ ᵔᴥᵔ ʔ,", a.Clients[a.clientID].Name, "(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧✧✧")
		a.actionsOfATM()

	case 2:
		var enterNumberOfLanguage int
		fmt.Println("1 - поменять на русский\n2 - поменять на английский")
		fmt.Scan(&enterNumberOfLanguage)
		a.ChangeLanguage(enterNumberOfLanguage)
		a.theFirstEntering()

	case 3:
		fmt.Println("До свидания")
	}
}

func (a *ATM) createNewClient(cardNum int) {
	var pin int
	var name string
	fmt.Println("Введите Имя и пинкод")
	fmt.Scan(&name, &pin)

	NewClient := Clients{
		amountOfWithdrawnMoney: 0,
		balance:                0,
		CardNumber:             cardNum,
		Name:                   name,
		PinCode:                pin,
	}

	a.clientID++
	a.Clients[a.clientID] = NewClient
}

func (a *ATM) checkClientForExistance(cardNum int) bool {
	for i, v := range a.Clients {
		if v.CardNumber == cardNum {
			a.clientID = i
			return true
		}
	}
	return false
}

func getBack() { // функциональность кнопки назад
	fmt.Println("1 - Назад")
	var num int
	fmt.Scan(&num)
	if num == 1 {
		return
	}
}

func initializeATM() ATM {
	atm := ATM{
		clientID:        0,
		language:        languageRussian,
		banknoteToCount: make(map[int]int),
		Clients:         make(map[int]Clients),
	}
	atm.banknoteToCount = map[int]int{100: 50, 500: 50, 1000: 50, 2000: 50, 5000: 50}

	for banknote, count := range atm.banknoteToCount {
		atm.balance += banknote * count
	}
	return atm
}

func main() {
	atm := initializeATM()
	atm.theFirstEntering()
}
