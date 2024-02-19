package jane

var Errors = map[string]string{
	`file_not_jane`:                `this is not jane source file: `,
	`invalid_token`:                `undenfied code content`,
	`invalid_syntax`:               `invalid syntax`,
	`exist_name`:                   `name is already exist`,
	`brace_not_closed`:             `brace is opened but not closed`,
	`function_body_not_exits`:      `function body is not declare`,
	`parameters_not_supported`:     `function is not support parameters`,
	`not_support_expression`:       `expression is not supports yet`,
	`missing_return`:               `missing return at end of function`,
	`invalid_numeric_range`:        `arithmetic value overflow`,
	`incompatible_value`:           `incompatible value with type`,
	`operator_overflow`:            `operator overflow`,
	`invalid_operator`:             `invalid operator`,
	`invalid_data_types`:           `data types are not compatible`,
	`operator_notfor_string`:       `this operator is not defined for string types`,
	`operator_notfor_booleans`:     `this operator is not defined for boolean types`,
	`operator_notfor_uint_and_int`: `this operator is not defined for uint and int types`,
}
