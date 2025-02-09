CREATE TABLE `evaluations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `departement_id` int(11) NOT NULL,
  `title` varchar(225) NOT NULL,
  `status` varchar(50) NOT NULL,
  `cc` varchar(225) DEFAULT NULL,
  `deadline_at` DATETIME DEFAULT NULL,
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


CREATE TABLE `evaluation_employees` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `employee_id` int(11) NOT NULL,
  `type` varchar(50) NOT NULL,
  `avg` float DEFAULT NULL,
  `email_sent` date DEFAULT NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `evaluation_answer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `evaluation_id` int(11) NOT NULL,
  `evaluation_employee_id` int(11) NOT NULL,
  `question_id` int(11) NOT NULL,
  `answer` text DEFAULT NULL,
  `final_point` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);