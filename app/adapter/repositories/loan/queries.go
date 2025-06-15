package loan

var (
	createLoan = `INSERT INTO loan (user_id, contract_number, otr_amount, principal_amount, asset_name, loan_type_id, limit_type_id, status, start_date, interest_rate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	createLoanPayment = `INSERT INTO loan_payment (loan_id, amount, date, channel) VALUES(?, ?, ?)`

	getLoanByContractNumber = `SELECT l.id, l.user_id, l.contract_number, l.otr_amount, l.principal_amount, l.asset_name, lot.name , lit.amount, lit.term , l.status, l.start_date, l.interest_rate FROM loan l JOIN limit_type lit ON l.limit_type_id = lit.id JOIN loan_type lot ON  l.loan_type_id  = lot.id WHERE l.user_id = ? AND l.contract_number = ?`

	getLoanPaymentsByLoanID = `SELECT amount, date, channel FROM loan_payment WHERE loan_id = ?`

	getLoanPaymentsByContractNumber = `WITH lpymnt_id AS (
		SELECT id FROM loan WHERE user_id = ? AND contract_number = ?
	)
	SELECT amount, date, channel FROM loan_payment WHERE loan_id IN (SELECT id FROM lpymnt_id)`
)
