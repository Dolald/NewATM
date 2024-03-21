package main

import (
	"fmt"
	"strconv"
	// печатание чека, давайте сохраним природу ?
)

type language string

var (
	languageRussian        language = "russian"
	languageEnglish        language = "english"
	banknote100            int      = 100
	banknote500            int      = 500
	banknote1000           int      = 1000
	banknote2000           int      = 2000
	banknote5000           int      = 5000
	sortedBanknotesStorage          = []int{5000, 1000, 2000, 500, 100}
	notSuggestionBanknote  int      = 1
)

type Client struct {
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
	Clients         map[int]Client
	language        language
}

func (a *ATM) changeLanguage(enterNumberOfLanguage int) {
	switch enterNumberOfLanguage {
	case 1:
		a.language = languageRussian
	case 2:
		a.language = languageEnglish
	}
}

func (a *ATM) checkBalance() {
	fmt.Println("Ваш баланс:", a.Clients[a.clientID].balance)
}

func (a *ATM) deposit(deposit int) {
	tempDeposit, estimateWithdrawanBanknoteAmount := deposit, 0

	for _, banknote := range sortedBanknotesStorage {
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

func (a *ATM) checkPinForWithdraw() bool {
	var pinCode int
	fmt.Println("Введите пинкод")
	fmt.Scan(&pinCode)

	if pinCode != a.Clients[a.clientID].PinCode {
		fmt.Println("Не правильный пинкод")
		return false
	}
	return true
}

func (a *ATM) displayAvalibleExchange(withdrawal, сhoosingExchange int) int {
	var banknotesForDisplay string
	if withdrawal >= 500 && withdrawal <= 120000 {
		fmt.Println("Купюры какого размена желаете ?\n1 - Всё равно")

		for banknote := range a.banknoteToCount {
			if withdrawal >= banknote && withdrawal < 6000 { // 6000, т.к. это максимальная сумма по минимальным купюром в данном кейсе
				banknotesForDisplay += strconv.Itoa(banknote) + " "
			} else if withdrawal <= banknote*60 && withdrawal >= 6000 {
				banknotesForDisplay += strconv.Itoa(banknote) + " "
			}
		}
		fmt.Println(banknotesForDisplay)
		fmt.Scan(&сhoosingExchange)
	} else {
		сhoosingExchange = 5000
	}
	return сhoosingExchange
}

func (a *ATM) isWithdrowalAmoutCorrect(withdrawal, ATMAmountForCheck int) bool {
	switch {
	case withdrawal%banknote100 != 0:
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

	if _, ex := a.banknoteToCount[banknote100]; !ex { // if нет 100, идёт проверка вводимого числа на кратность 100, но не 500 или вывод < 500
		if withdrawal%banknote100 == 0 && withdrawal%banknote500 != 0 || withdrawal < banknote500 {
			fmt.Println("Введите другую сумму")
			return false
		}
	}
	return true
}

func (a *ATM) withdrawOneBanknote(chousenBanknote int) {
	totalAmountOneBanknoteMoney := chousenBanknote * a.banknoteToCount[chousenBanknote]
	if !a.isWithdrowalAmoutCorrect(chousenBanknote, totalAmountOneBanknoteMoney) {
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

func (a *ATM) getBancknotes() string {
	var displayBanknotes string
	for banknote := range a.banknoteToCount {
		displayBanknotes += strconv.Itoa(banknote) + " "
	}
	return displayBanknotes
}

func (a *ATM) withdraw() {
	var withdrawal, chousenBanknote, сhoosingExchange, ATMAmountForCheck int

	avalibleBanknotes := a.getBancknotes()
	fmt.Println("Выберите сумму для снятия\n1 - ввести другую сумму\n", avalibleBanknotes)
	fmt.Scan(&chousenBanknote)

	if !a.checkPinForWithdraw() {
		return
	}

	if chousenBanknote != notSuggestionBanknote {
		a.withdrawOneBanknote(chousenBanknote)
		return
	}

	fmt.Println("Введите сумму для снятия наличных\nВведите сумму до 300000\nДоступны номиналы купюр:", avalibleBanknotes)
	fmt.Scan(&withdrawal)

	сhoosingExchange = a.displayAvalibleExchange(withdrawal, сhoosingExchange)

	for banknote, ammount := range a.banknoteToCount {
		if banknote <= сhoosingExchange {
			ATMAmountForCheck += banknote * ammount
		}
	}

	if !(a.isWithdrowalAmoutCorrect(withdrawal, ATMAmountForCheck) && a.isWithdrowalAmoutBanknotesCorrect(сhoosingExchange, withdrawal)) {
		return
	}

	a.withdrawMoneyFromATM(сhoosingExchange, withdrawal)
	a.updateClientOptions(withdrawal)
	a.balance -= withdrawal
	fmt.Println("Выдача купюр. . .\nВы успешно сняли", withdrawal)
}

func (a *ATM) isWithdrowalAmoutBanknotesCorrect(сhoosingExchange, withdrawal int) bool {
	var countWithdrawalBanknotes, WithdrawalBanknotes int
	for _, banknote := range sortedBanknotesStorage {
		if withdrawal != 0 && banknote <= сhoosingExchange {

			WithdrawalBanknotes = a.banknoteToCount[banknote]
			WithdrawalBanknotes -= withdrawal / banknote
			countWithdrawalBanknotes += a.banknoteToCount[banknote] - WithdrawalBanknotes

			if a.banknoteToCount[banknote] < 0 {
				a.banknoteToCount[banknote] = 0
			}

			if withdrawal/banknote != 0 {
				withdrawal -= (a.banknoteToCount[banknote] - WithdrawalBanknotes) * banknote
			}
		}
	}
	if countWithdrawalBanknotes > 60 {
		fmt.Println("банкомат не может выдать такое количество купюр")
		return false
	}
	return true
}

func (a *ATM) updateClientOptions(estimatedWithdrawal int) {
	if v, ok := a.Clients[a.clientID]; ok {
		v.amountOfWithdrawnMoney += estimatedWithdrawal
		v.balance -= estimatedWithdrawal
		a.Clients[a.clientID] = v
	}
}

func (a *ATM) doesAllBanknotesExist() bool {
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

func (a *ATM) withdrawMoneyFromATM(сhoosingExchange, withdrawal int) {
	var tempEstimatedWithdrawalAmount int
	for _, banknote := range sortedBanknotesStorage {
		if withdrawal != 0 && banknote <= сhoosingExchange {

			tempEstimatedWithdrawalAmount = a.banknoteToCount[banknote]
			a.banknoteToCount[banknote] -= withdrawal / banknote

			if a.banknoteToCount[banknote] < 0 {
				a.banknoteToCount[banknote] = 0
			}

			if withdrawal/banknote != 0 {
				withdrawal -= (tempEstimatedWithdrawalAmount - a.banknoteToCount[banknote]) * banknote
			}
		}
	}
}

func (a *ATM) actionsOfATM() {
	fmt.Println("1 - снять наличные \n2 - пополнить карту \n3 - посмотреть баланс \n4 - поменять язык \n5 - вытащить карту")
	var enterNum int
	fmt.Scan(&enterNum)

	switch enterNum {
	case 1:
		a.withdraw()
		if !a.doesAllBanknotesExist() {
			return
		}
		getBack()
		a.actionsOfATM()

	case 2:
		var deposit int
		fmt.Println("Введите сумму для пополнения")
		fmt.Scan(&deposit)
		if deposit%banknote100 != 0 {
			fmt.Println("Банкомата не принимает такие суммы, выберите другую сумму")
			return
		}
		a.deposit(deposit)
		getBack()
		a.actionsOfATM()

	case 3:
		a.checkBalance()
		getBack()
		a.actionsOfATM()

	case 4:
		var enterNumberOfLanguage int
		fmt.Println("1 - поменять на русский\n2 - поменять на английский")
		fmt.Scan(&enterNumberOfLanguage)
		a.changeLanguage(enterNumberOfLanguage)
		getBack()
		a.actionsOfATM()

	case 5:
		fmt.Println("Заберите карту")
		a.theFirstEntering()
	}
}

func (a *ATM) isPinСodeCorrect() bool {
	var pin, countWrongPin int
	for {
		fmt.Println("Введите пинкод")
		fmt.Scan(&pin)

		if a.Clients[a.clientID].PinCode == pin {
			break
		}
		countWrongPin++
		if countWrongPin == 3 {
			fmt.Println("Ваша карта забанена")
			return false
		}
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

		if a.doesClientExist(cardNum) {
			if !a.isPinСodeCorrect() {
				return
			}
			a.identifyTheСlient(cardNum)
		} else {
			a.createNewClient(cardNum)
		}

		fmt.Println("Добро пожаловать ʕ ᵔᴥᵔ ʔ,", a.Clients[a.clientID].Name)
		a.actionsOfATM()

	case 2:
		var enterNumberOfLanguage int
		fmt.Println("1 - поменять на русский\n2 - поменять на английский")
		fmt.Scan(&enterNumberOfLanguage)
		a.changeLanguage(enterNumberOfLanguage)
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

	NewClient := Client{
		amountOfWithdrawnMoney: 0,
		balance:                0,
		CardNumber:             cardNum,
		Name:                   name,
		PinCode:                pin,
	}

	a.clientID++
	a.Clients[a.clientID] = NewClient
}

func (a *ATM) identifyTheСlient(cardNum int) {
	for i, v := range a.Clients {
		if v.CardNumber == cardNum {
			a.clientID = i
		}
	}
}

func (a *ATM) doesClientExist(cardNum int) bool {
	for _, v := range a.Clients {
		if v.CardNumber == cardNum {
			return true
		}
	}
	return false
}

func getBack() { // функциональность кнопки назад
	var num int
	fmt.Println("1 - Назад")
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
		Clients:         make(map[int]Client),
	}
	atm.banknoteToCount = map[int]int{banknote100: 50, banknote500: 50, banknote1000: 50, banknote2000: 50, banknote5000: 50}

	for banknote, count := range atm.banknoteToCount {
		atm.balance += banknote * count
	}
	return atm
}

func main() {
	atm := initializeATM()
	atm.theFirstEntering()
}
