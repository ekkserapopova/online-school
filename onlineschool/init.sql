INSERT INTO Tasks (name, description, module_id)
VALUES
    ('Факторизация чисел', 'Напишите функцию на Python, которая находит все простые множители заданного числа и возвращает их в виде списка.', 3),
    ('Реализация сложения матриц', 'Создайте функцию на Python, которая выполняет сложение двух матриц одинаковой размерности. Функция должна принимать два двумерных массива и возвращать их сумму.',3),
    ('Поиск максимального элемента', 'Создайте программу на Python для нахождения максимального элемента в списке чисел.', 2)


-- Questions about conditional operators and loops for module_id=2
INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Какой условный оператор используется для проверки одного условия в большинстве языков программирования?', 1, 2, 1, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Что происходит, когда в условии оператора if указано выражение, возвращающее false?', null, 2, 1, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Какой цикл подходит для выполнения блока кода заранее известное количество раз?', null, 2, 1, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Что такое бесконечный цикл и как его можно создать?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Какой оператор используется для прерывания выполнения текущей итерации цикла и перехода к следующей?', null, 2, 1, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Для чего используется оператор switch/case?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('В чем разница между циклами while и do-while?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Какой оператор позволяет полностью прекратить выполнение цикла?', null, 2, 1, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Что такое вложенные условные операторы? Приведите пример использования.', null, 2, 3, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Как можно объединить несколько условий в операторе if?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Что такое тернарный оператор условия? В каких случаях его рекомендуется использовать?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Какие проблемы могут возникнуть при использовании вложенных циклов?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Чем отличается условный оператор if от тернарного оператора?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Как работает цикл for-each (или foreach) и в каких случаях его удобно применять?', null, 2, 2, NOW(), NOW());

INSERT INTO questions (text, test_id, module_id, points, created_at, updated_at)
VALUES ('Что такое "короткое замыкание" (short-circuit) в условных выражениях и как оно работает?', null, 2, 3, NOW(), NOW());


-- Answers for Question 1: Какой условный оператор используется для проверки одного условия в большинстве языков программирования?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'if', true, 1, 1, NOW(), NOW()),
    (2, 'switch', false, 1, 1, NOW(), NOW()),
    (3, 'for', false, 1, 1, NOW(), NOW()),
    (4, 'while', false, 1, 1, NOW(), NOW());

-- Answers for Question 52: Что происходит, когда в условии оператора if указано выражение, возвращающее false?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Выполняется блок кода внутри if', false, 52, null, NOW(), NOW()),
    (2, 'Выполняется блок кода в else, если он существует', true, 52, NULL, NOW(), NOW()),
    (3, 'Программа завершается с ошибкой', false, 52, NULL, NOW(), NOW()),
    (4, 'Код внутри условного оператора игнорируется и выполнение продолжается со следующей строки после блока if', true, 2, NULL, NOW(), NOW());

-- Answers for Question 3: Какой цикл подходит для выполнения блока кода заранее известное количество раз?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'for', true, 53, NULL, NOW(), NOW()),
    (2, 'while', false, 53, NULL, NOW(), NOW()),
    (3, 'do-while', false, 53, NULL, NOW(), NOW()),
    (4, 'if-else', false, 53, NULL, NOW(), NOW());

-- Answers for Question 4: Что такое бесконечный цикл и как его можно создать?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Цикл, который никогда не завершается без внешнего вмешательства. Можно создать с помощью while(true) или for(;;)', true, 54, 61, NOW(), NOW()),
    (2, 'Цикл, который выполняется ровно 100 раз', false, 54, 61, NOW(), NOW()),
    (3, 'Цикл, который автоматически останавливается при переполнении памяти', false, 54, 61, NOW(), NOW()),
    (4, 'Цикл, выполняющийся только один раз', false, 54, 61, NOW(), NOW());

-- Answers for Question 5: Какой оператор используется для прерывания выполнения текущей итерации цикла и перехода к следующей?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'break', false, 55, NULL, NOW(), NOW()),
    (2, 'continue', true, 55, NULL, NOW(), NOW()),
    (3, 'return', false, 55, NULL, NOW(), NOW()),
    (4, 'next', false, 55, NULL, NOW(), NOW());

-- Answers for Question 6: Для чего используется оператор switch/case?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Для организации циклов', false, 56, NULL, NOW(), NOW()),
    (2, 'Для проверки множества возможных значений одной переменной', true, 56, NULL, NOW(), NOW()),
    (3, 'Для подключения внешних библиотек', false, 56, NULL, NOW(), NOW()),
    (4, 'Для обработки исключений', false, 56, NULL, NOW(), NOW());

-- Answers for Question 7: В чем разница между циклами while и do-while?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Они абсолютно идентичны', false, 57, NULL, NOW(), NOW()),
    (2, 'while проверяет условие до выполнения тела цикла, а do-while - после', true, 57, NULL, NOW(), NOW()),
    (3, 'do-while выполняется бесконечно, а while останавливается автоматически', false, 57, NULL, NOW(), NOW()),
    (4, 'while работает быстрее, чем do-while', false, 57, NULL, NOW(), NOW());

-- Answers for Question 8: Какой оператор позволяет полностью прекратить выполнение цикла?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'break', true, 58, NULL, NOW(), NOW()),
    (2, 'continue', false, 58, NULL, NOW(), NOW()),
    (3, 'exit', false, 58, NULL, NOW(), NOW()),
    (4, 'stop', false, 58, NULL, NOW(), NOW());

-- Answers for Question 9: Что такое вложенные условные операторы? Приведите пример использования.
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Это использование условных операторов внутри других условных операторов. Например: if (возраст >= 18) { if (есть_права) { разрешить_вождение(); } }', true, 59, NULL, NOW(), NOW()),
    (2, 'Это использование нескольких условий в одном операторе if', false, 59, NULL, NOW(), NOW()),
    (3, 'Это использование оператора switch внутри цикла for', false, 59, NULL, NOW(), NOW()),
    (4, 'Это условия, которые всегда возвращают true', false, 59, NULL, NOW(), NOW());

-- Answers for Question 10: Как можно объединить несколько условий в операторе if?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Используя операторы && (И) и || (ИЛИ)', true, 60, NULL, NOW(), NOW()),
    (2, 'Написав несколько if подряд', false, 60, NULL, NOW(), NOW()),
    (3, 'Используя оператор switch', false, 60, NULL, NOW(), NOW()),
    (4, 'Условия нельзя объединять', false, 60, NULL, NOW(), NOW());

-- Answers for Question 11: Что такое тернарный оператор условия? В каких случаях его рекомендуется использовать?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Оператор, который позволяет написать условие в одну строку вида условие ? выражение1 : выражение2. Рекомендуется использовать для простых условий', true, 61, 61, NOW(), NOW()),
    (2, 'Оператор, который может содержать три условия одновременно', false, 61, 61, NOW(), NOW()),
    (3, 'Оператор для работы с тремя переменными', false, 61, 61, NOW(), NOW()),
    (4, 'Специальная форма цикла for с тремя параметрами', false, 61, 61, NOW(), NOW());

-- Answers for Question 12: Какие проблемы могут возникнуть при использовании вложенных циклов?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Низкая производительность из-за экспоненциального роста количества операций', true, 62, NULL, NOW(), NOW()),
    (2, 'Ошибки компиляции', false, 62, NULL, NOW(), NOW()),
    (3, 'Вложенные циклы не поддерживаются в современных языках программирования', false, 62, NULL, NOW(), NOW()),
    (4, 'Автоматическое завершение программы', false, 62, NULL, NOW(), NOW());

-- Answers for Question 13: Чем отличается условный оператор if от тернарного оператора?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Тернарный оператор всегда возвращает значение и записывается в одну строку, if может не возвращать значение и занимает больше строк', true, 63, NULL, NOW(), NOW()),
    (2, 'Тернарный оператор работает быстрее, чем if', false, 63, NULL, NOW(), NOW()),
    (3, 'If может содержать только одно условие, а тернарный - множество', false, 63, NULL, NOW(), NOW()),
    (4, 'Они полностью идентичны по функциональности', false, 63, NULL, NOW(), NOW());

-- Answers for Question 14: Как работает цикл for-each (или foreach) и в каких случаях его удобно применять?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'For-each перебирает все элементы коллекции без необходимости работы с индексами. Удобен для простого перебора элементов массива или коллекции', true, 64, 61, NOW(), NOW()),
    (2, 'For-each выполняет цикл ровно 100 раз', false, 64, 61, NOW(), NOW()),
    (3, 'For-each - это устаревший тип цикла, который не рекомендуется использовать', false, 64, 61, NOW(), NOW()),
    (4, 'For-each выполняется быстрее, чем обычный цикл for', false, 64, 61, NOW(), NOW());

-- Answers for Question 15: Что такое "короткое замыкание" (short-circuit) в условных выражениях и как оно работает?
INSERT INTO answer_variants (number_id, text, is_right, question_id, test_id, created_at, updated_at)
VALUES
    (1, 'Это принцип оптимизации, при котором вычисление логического выражения прерывается, как только становится известен результат. Например, в выражении (A && B), если A = false, B не вычисляется', true, 65, 61, NOW(), NOW()),
    (2, 'Это ошибка в программе, когда условие никогда не выполняется', false, 65, 61, NOW(), NOW()),
    (3, 'Это особый тип цикла, который завершается досрочно', false, 65, 61, NOW(), NOW()),
    (4, 'Это метод обработки исключений в условных выражениях', false, 65, 61, NOW(), NOW());
\


select avg(score) from student_tasks where student_id = 1 and task_id = 3