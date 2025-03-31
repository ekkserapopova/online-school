-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    surname VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    birth DATE,
    photo VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    overview TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Таблица языков
CREATE TABLE IF NOT EXISTS languages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);



-- Таблица курсов
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    difficulty INTEGER NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    teacher_id INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE
);


-- Таблица для связи многие-ко-многим между языками и курсами
CREATE TABLE IF NOT EXISTS language_courses (
    language_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    PRIMARY KEY (language_id, course_id),
    FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Таблица для связи многие-ко-многим между студентами и курсами
CREATE TABLE IF NOT EXISTS students_courses (
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    PRIMARY KEY (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Таблица расписания
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица для связи многие-ко-многим между расписаниями и курсами
CREATE TABLE IF NOT EXISTS schedules_courses (
    schedule_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    PRIMARY KEY (schedule_id, course_id),
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Таблица уроков
CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    start TIMESTAMP WITH TIME ZONE NOT NULL,
    "end" TIMESTAMP WITH TIME ZONE NOT NULL,
    course_id INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Таблица материалов
CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    file VARCHAR(255),
    lesson_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);

-- Таблица домашних заданий
CREATE TABLE IF NOT EXISTS homeworks (
    id SERIAL PRIMARY KEY,
    deadline TIMESTAMP WITH TIME ZONE,
    file VARCHAR(255),
    result INTEGER,
    comment TEXT,
    implementation_status VARCHAR(50),
    lesson_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица тестов
CREATE TABLE IF NOT EXISTS tests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    lesson_id INTEGER,
    course_id INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE SET NULL,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Таблица вопросов
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    text TEXT,
    answer VARCHAR(255),
    test_id INTEGER NOT NULL,
    points INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
);

-- Таблица ответов
CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    student_answer TEXT,
    result BOOLEAN,
    points_earned INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица отзывов
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    text TEXT,
    assessment INTEGER NOT NULL CHECK (assessment BETWEEN 1 AND 5),
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    amount INTEGER,
    status VARCHAR(50) NOT NULL,
    date TIMESTAMP WITH TIME ZONE,
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);


UPDATE users SET is_active = TRUE WHERE is_active IS NULL;
UPDATE users SET is_admin = FALSE WHERE is_admin IS NULL;
UPDATE courses SET is_active = TRUE WHERE is_active IS NULL;
UPDATE lessons SET is_active = TRUE WHERE is_active IS NULL;
UPDATE tests SET is_active = TRUE WHERE is_active IS NULL;
UPDATE reviews SET is_published = FALSE WHERE is_published IS NULL;
UPDATE homeworks SET status = TRUE WHERE status IS NULL;

-- Триггер для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    t text;
BEGIN
    FOR t IN 
        SELECT table_name FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        AND table_name IN ('users', 'languages', 'courses', 'lessons', 
                          'materials', 'homeworks', 'tests', 'questions', 
                          'answers', 'reviews', 'payments', 'schedules')
    LOOP
        EXECUTE format('
            DROP TRIGGER IF EXISTS set_%I_updated_at ON %I;
            CREATE TRIGGER set_%I_updated_at
            BEFORE UPDATE ON %I
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
        ', t, t, t, t);
    END LOOP;
END;
$$;



insert into languages (name) VALUES 
('Python'),
('JavaScript'),
('Java'),
('C++'),
('Go'),
('PHP'),
('Ruby'),
('Swift'),
('Kotlin'),
('Dart');

INSERT INTO courses (name, description, difficulty, price, teacher_id, is_active) VALUES 
('Основы Python для начинающих', 'Полный курс с нуля для тех, кто хочет освоить самый популярный язык программирования. Изучите синтаксис, структуры данных, работу с файлами и основы ООП.', 1, 0, 1, true),
('JavaScript и React для веб-разработчиков', 'Интенсивный курс по созданию современных веб-приложений на React. От основ JavaScript до продвинутых паттернов React и управления состоянием.', 3, 0, 1, true),
('Алгоритмы и структуры данных на Go', 'Погрузитесь в мир эффективных алгоритмов и структур данных с использованием Go. Курс для тех, кто хочет писать оптимальный код и пройти технические собеседования.', 4, 0, 1, true),
('Разработка мобильных приложений на Flutter', 'Научитесь создавать кроссплатформенные мобильные приложения на Flutter и Dart. От основ до публикации в App Store и Google Play.', 3, 0, 1, true),
('DevOps и CI/CD с Docker и Kubernetes', 'Полное руководство по автоматизации процессов разработки, тестирования и развертывания. Освойте Docker, Kubernetes, GitHub Actions и другие современные инструменты.', 5, 0, 1, true);
