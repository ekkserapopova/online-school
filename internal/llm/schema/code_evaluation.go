package schema

type CodeEvaluation struct {
	Requirements     string  `json:"requirements" jsonschema_description:"Анализ соответствия явным требованиям задания, вывод оценки (оценка 0-5)"`
	Implementation   string  `json:"implementation" jsonschema_description:"Анализ корректности реализации, вывод оценки (оценка 0-5)"`
	BoundaryHandling string  `json:"boundary_handling" jsonschema_description:"Анализ обработки граничных случаев и ошибок, вывод оценки (оценка 0-5)"`
	Optimization     string  `json:"optimization" jsonschema_description:"Анализ эффективности алгоритма и оптимизации, вывод оценки (оценка 0-5)"`
	Formula          string  `json:"formula" jsonschema_description:"Формула для подсчета итоговой оценки с подстановкой всех оценк, в численном варианте"`
	Score            float32 `json:"score" jsonschema_description:"Итоговая числовая оценка кода от 0 до 5"`
	Conclusion       string  `json:"conclusion" jsonschema_description:"Итоговое заключение с обоснованием оценки и рекомендациями"`
}
