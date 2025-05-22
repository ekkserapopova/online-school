package schema

type CodeEvaluation struct {
	Requirements      string `json:"requirements" jsonschema_description:"Анализ соответствия явным требованиям задания, вывод оценки (оценка 0-5)"`
	RequirementsScore int    `json:"requirements_score" jsonschema_description:"Оценка по критерию 'Корректность реализации' от 0 до 5"`

	Implementation      string `json:"implementation" jsonschema_description:"Анализ корректности реализации, вывод оценки (оценка 0-5)"`
	ImplementationScore int    `json:"implementation_score" jsonschema_description:"Оценка по критерию 'Качество кода' от 0 до 5"`

	BoundaryHandling      string `json:"boundary_handling" jsonschema_description:"Анализ обработки граничных случаев и ошибок, вывод оценки (оценка 0-5)"`
	BoundaryHandlingScore int    `json:"boundary_handling_score" jsonschema_description:"Оценка по критерию 'Обработка краевых случаев' от 0 до 5"`

	Optimization      string `json:"optimization" jsonschema_description:"Анализ эффективности алгоритма и оптимизации, вывод оценки (оценка 0-5)"`
	OptimizationScore int    `json:"optimization_score" jsonschema_description:"Оценка по критерию 'Эффективность алгоритма' от 0 до 5"`

	//Formula          string  `json:"formula" jsonschema_description:"Формула для подсчета итоговой оценки с подстановкой всех оценк, в численном варианте"`
	//Score            float32 `json:"score" jsonschema_description:"Итоговая числовая оценка кода от 0 до 5"`
	Conclusion      string `json:"conclusion" jsonschema_description:"Итоговое заключение с рекомендациями"`
	CodeWithComment string `json:"code_with_comment" jsonschema_description:"Исходный код ученика с добавлением комментария в строку с ошибкой"`
}
