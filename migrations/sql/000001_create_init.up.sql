CREATE TABLE `evaluations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `departement_id` int(11) NOT NULL,
  `title` varchar(225) NOT NULL,
  `status` varchar(50) NOT NULL,
  `deadline_at` varchar(50) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
);


CREATE TABLE `questions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `title` varchar(225) NULL,
  `question` text NOT NULL,
  `type` varchar(50) NOT NULL,
  `competency_type` varchar(50) NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `evaluated_employees` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `employee_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `evaluator_employees` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `evaluated_employee_id` int(11) NOT NULL,
  `employee_id` int(11) NOT NULL,
  `total_functional` float DEFAULT 0,
  `total_personal` float DEFAULT 0,
  `total_avg` float DEFAULT 0,
  `has_assessed` boolean DEFAULT false,
  `requires_assessment` boolean DEFAULT false,
  `email_sent` varchar(50) DEFAULT NULL,
  `cc` varchar(225) DEFAULT NULL,
  `status` varchar(50) DEFAULT 'pending',
  PRIMARY KEY (`id`)
);


CREATE TABLE `evaluation_answers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `evaluator_employee_id` int(11) NOT NULL,
  `question_id` int(11) NOT NULL,
  `answer` text DEFAULT NULL,
  `final_point` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE evaluated_employees ADD total_avg FLOAT DEFAULT 0;


