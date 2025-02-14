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
  `question` text NOT NULL,
  `type` varchar(50) NOT NULL,
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
  `avg` float DEFAULT NULL,
  `email_sent` varchar(50) DEFAULT NULL,
  `cc` varchar(225) DEFAULT NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `evaluation_answers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `evaluation_employee_id` int(11) NOT NULL,
  `question_id` int(11) NOT NULL,
  `answer` text DEFAULT NULL,
  `final_point` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);

