package repo

import (
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"math"
	"math/rand"
	"onlineschool/internal/models"
)

func (r *Repo) GetTest(testID int) (models.Test, error) {
	var test models.Test
	err := r.db.Preload("Questions.Answers").Where("id = ?", testID).First(&test).Error
	if err != nil {
		return test, err
	}

	return test, nil
}

func (r *Repo) GetTests(courseID int) ([]models.Test, error) {
	var tests []models.Test
	err := r.db.Preload("Questions").Where("module_id = ?", courseID).Find(&tests).Error
	if err != nil {
		return tests, err
	}

	return tests, nil
}

func (r *Repo) GetQuestionsForTest(testID int) ([]models.Question, error) {
	questions := []models.Question{}

	err := r.db.Preload("Answers").Where("test_id = ?", testID).Find(&questions).Error

	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *Repo) GetQuestion(questionID int) (models.Question, error) {
	qusetion := models.Question{}
	err := r.db.Preload("Answers").Where("id = ?", questionID).First(&qusetion).Error
	if err != nil {
		return qusetion, err
	}

	return qusetion, nil
}

func (r *Repo) ExistsStudentAnswer(studentID, questionID int) (bool, error) {
	studentAnswer := models.StudentAnswer{}

	err := r.db.Where("student_id = ? AND question_id = ?", studentID, questionID).First(&studentAnswer).Error

	log.Println(studentAnswer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // запись не найдена, это не ошибка
		}
		return false, err // другая ошибка базы данных
	}

	return true, nil

}

func (r *Repo) AddAnswerByStudent(studentAnswer models.StudentAnswer) (models.StudentAnswer, error) {
	//TODO: ученик пролходит тест второй раз
	questionID := studentAnswer.QuestionID
	existsStudentAnswer, err := r.ExistsStudentAnswer(studentAnswer.StudentID, questionID)

	log.Printf("Student answer was got: %+v", existsStudentAnswer)

	if err != nil {
		return studentAnswer, err
	}

	question, err := r.GetQuestion(questionID)
	if err != nil {
		return studentAnswer, err
	}

	studentAnswer.TestID = question.TestID

	rightAnswers := []int{}

	for _, answer := range question.Answers {
		if answer.IsRight {
			rightAnswers = append(rightAnswers, answer.ID)
		}
	}

	selectedAnswersInt := int64ArrayToIntSlice(studentAnswer.SelectedAnswerIDs)

	// Теперь сравниваем слайсы одинакового типа
	studentAnswer.Result = compareSlices(selectedAnswersInt, rightAnswers)
	log.Printf("Student answer was got: %+v", selectedAnswersInt)
	log.Printf("Right answer was got: %+v", rightAnswers)

	// Начисление баллов
	if studentAnswer.Result {
		studentAnswer.PointsEarned = question.Points
	} else {
		studentAnswer.PointsEarned = 0
	}

	err = r.db.Save(&studentAnswer).Error
	// Используем транзакцию для гарантированного обновления
	//tx := r.db.Begin()
	if err != nil {
		//tx.Rollback()
		log.Printf("Error occurred: %v", err)
		return studentAnswer, err
	}

	//tx.Commit()
	log.Printf("Student answer was successfully created or updated")

	var completedTest models.CompletedTest

	err = r.db.Where("student_id = ? and test_id = ?", studentAnswer.StudentID, studentAnswer.TestID).
		First(&completedTest).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			completedTest.TestID = studentAnswer.TestID
			completedTest.StudentID = studentAnswer.StudentID
			completedTest.Status = "in progress"
			completedTest.Points += studentAnswer.PointsEarned
			err = r.db.Create(&completedTest).Error
			if err != nil {
				return studentAnswer, err
			}
		}
		return studentAnswer, err
	}

	completedTest.Points += studentAnswer.PointsEarned

	err = r.db.Save(&completedTest).Error
	if err != nil {
		return studentAnswer, err
	}

	return studentAnswer, nil
}
func (r *Repo) GetPointsOfTest(testID, studentID int) (int, error) {

	var totalPoints int
	err := r.db.Model(&models.StudentAnswer{}).
		Where("test_id = ? AND student_id = ?", testID, studentID).
		Select("COALESCE(sum(points_earned), 0) as total_points").
		Scan(&totalPoints).Error

	if err != nil {
		return 0, err
	}

	return totalPoints, nil
}

func (r *Repo) GetStudentAnswers(testID int, studentID int) ([]models.StudentAnswer, error) {
	//questions, err := r.GetQuestionsForTest(testID)
	//if err != nil {
	//	return nil, err
	//}

	studentAnswers := []models.StudentAnswer{}

	err := r.db.
		Where("test_id = ? and student_id = ?", testID, studentID).Find(&studentAnswers).Error

	if err != nil {
		return nil, err
	}

	return studentAnswers, nil
}

func (r *Repo) GetRightAnswers(testID int) ([]models.AnswerResponse, error) {
	//TODO: preload answer's info
	var rightAnswers []models.AnswerResponse

	err := r.db.
		Model(&models.AnswerVariant{}).
		Where("test_id = ? AND is_right = ?", testID, true).
		Find(&rightAnswers).Error

	if err != nil {
		return nil, err
	}

	return rightAnswers, nil
}

func (r *Repo) GetTestResultForStudent(testID, studentID int) (struct {
	Test         models.Test
	RightAnswers []models.AnswerVariant
	Points       int
}, error) {

	type response struct {
		Test         models.Test
		RightAnswers []models.AnswerVariant
		Points       int
	}

	rightAnswers := []models.AnswerVariant{}

	test, err := r.GetTest(testID)
	if err != nil {
		return response{test, rightAnswers, 0}, err
	}

	questions := test.Questions

	for _, question := range questions {
		for _, answer := range question.Answers {
			if answer.IsRight == true {
				rightAnswers = append(rightAnswers, answer)
			}
		}
	}

	points, err := r.GetPointsOfTest(testID, studentID)
	if err != nil {
		return response{test, rightAnswers, 0}, err
	}

	return response{test, rightAnswers, points}, nil
}

func compareSlices(selected, correct []int) bool {
	// Если количество выбранных ответов не совпадает с количеством правильных, сразу возвращаем false
	if len(selected) != len(correct) {
		return false
	}

	// Создаем множество правильных ответов для быстрого поиска
	correctSet := make(map[int]bool)
	for _, v := range correct {
		correctSet[v] = true
	}

	// Проверяем, что каждый выбранный ответ есть среди правильных
	for _, v := range selected {
		if !correctSet[v] {
			return false
		}
	}

	// Все выбранные ответы правильные и их количество совпадает с количеством правильных ответов
	return true
}

func int64ArrayToIntSlice(arr pq.Int64Array) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}

func (r *Repo) FinishTest(testID int, studentID int) error {
	var completedTest models.CompletedTest

	err := r.db.
		//Preload("Questions.StudentAnswers").
		Where("test_id = ? AND student_id = ?", testID, studentID).
		First(&completedTest).Error

	if err != nil {
		log.Println("Student answer was not found")
		return err
	}

	log.Println("Student answer was successfully finished or updated")

	completedTest.Status = "completed"

	return r.db.Save(&completedTest).Error
}

func (r *Repo) GetCompletedTest(testID int, studentID int) (models.CompletedTest, error) {
	var completedTest models.CompletedTest
	err := r.db.
		Where("test_id = ? AND student_id = ?", testID, studentID).
		First(&completedTest).Error

	if err != nil {
		log.Println("Student answer was not found")
		return completedTest, err
	}

	if completedTest.Status != "completed" {
		log.Println("Test is not completed")
		return completedTest, errors.New("Test is not completed")
	}

	return completedTest, nil
}

func (r *Repo) RandomGenerateTest(testID int) (models.Test, error) {
	questions := []models.Question{}
	test, err := r.GetTest(testID)
	if err != nil {
		return test, err
	}

	// Получаем доступные вопросы без привязки к тесту
	err = r.db.
		Where("module_id = ? AND test_id IS NULL", test.ModuleID).
		Find(&questions).Error
	if err != nil {
		return test, err
	}

	if len(questions) < int(math.Abs(float64(test.CountQuestions))) {
		test.CountQuestions = len(questions)
	}

	// Перемешиваем вопросы
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	questionsForNewTest := questions[:int(math.Abs(float64(test.CountQuestions)))]

	// Обновляем каждый вопрос: test_id и module_id = NULL
	for i := range questionsForNewTest {
		questionsForNewTest[i].TestID = &test.ID

		// Сохраняем test_id и сбрасываем module_id через прямой SQL
		err := r.db.Model(&questionsForNewTest[i]).
			Updates(map[string]interface{}{
				"test_id":   test.ID,
				"module_id": gorm.Expr("NULL"),
			}).Error
		if err != nil {
			return test, err
		}
	}

	return test, nil
}

func (r *Repo) CreateTest(test models.Test) (models.Test, error) {
	err := r.db.Create(&test).Error
	return test, err
}

func (r *Repo) DeleteTest(testID int) error {
	err := r.db.Delete(&models.Test{}, "id = ?", testID).Error
	return err
}

func (r *Repo) GetQuestionsForModule(moduleID int) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Preload("Answers").Where("module_id = ? and test_id is null", moduleID).Find(&questions).Error
	return questions, err
}

func (r *Repo) AddQuestion(question models.Question) (models.Question, error) {
	err := r.db.Create(&question).Error
	return question, err
}

func (r *Repo) AddAnswer(answer models.AnswerVariant) (models.AnswerVariant, error) {
	err := r.db.Create(&answer).Error
	return answer, err
}

func (r *Repo) UpdateTest(test models.Test) error {
	err := r.db.Save(&test).Error
	return err
}

func (r *Repo) CreateStudentTest(testID, studentID int) (models.CompletedTest, error) {
	var test models.CompletedTest
	test.TestID = &testID
	test.Status = "in progress"
	test.StudentID = studentID
	test.Points = 0
	err := r.db.Save(&test).Error
	return test, err
}

func (r *Repo) GetStudentsTests(testID int) ([]models.CompletedTestResponse, error) {
	var tests []models.CompletedTest
	err := r.db.Where("test_id = ?", testID).Find(&tests).Error
	var testsResult []models.CompletedTestResponse
	for _, test := range tests {
		var student models.User
		err = r.db.Where("id = ?", test.StudentID).First(&student).Error
		student_name := student.Name
		student_surname := student.Surname
		testsResult = append(testsResult, models.CompletedTestResponse{
			ID:             test.ID,
			TestID:         &test.ID,
			StudentID:      test.StudentID,
			Status:         test.Status,
			Points:         test.Points,
			CreatedAt:      test.CreatedAt,
			UpdatedAt:      test.UpdatedAt,
			StudentName:    student_name,
			StudentSurname: student_surname,
		})
	}
	return testsResult, err
}

func (r *Repo) DeleteAnswer(answerID int) error {
	err := r.db.Where("id = ?", answerID).Delete(&models.AnswerVariant{}).Error
	return err
}

func (r *Repo) DeleteQuestion(questionID int) error {
	err := r.db.Where("id = ?", questionID).Delete(&models.Question{}).Error
	return err
}

func (r *Repo) UpdateQuestion(question models.Question) error {
	err := r.db.Save(&question).Error
	return err
}

func (r *Repo) GetAnswerVariant(answerID int) (models.AnswerVariant, error) {
	var answer models.AnswerVariant
	err := r.db.Where("id = ?", answerID).First(&answer).Error
	return answer, err
}

func (r *Repo) UpdateAnswerVariant(answerVariant models.AnswerVariant) error {
	err := r.db.Save(&answerVariant).Error
	return err
}
