SELECT * from managers;

SELECT * from users;

SELECT * from courses;

SELECT * from users_tokens;

SELECT * from user_courses order by user_id;

SELECT EXISTS (SELECT course_name FROM user_courses WHERE user_id = 13 AND courses_id = 1);

update courses SET course_name = 'German course lvl 1' where id = 4;

SELECT is_admin from managers WHERE id = 8;



Manager Admin
{
    "name":     "Frank",
    "phone":    "+992900800664",
    "login":    "adminManager",
    "password": "adas12fasfasdf1asf"
}



SELECT (course_name, price) FROM courses WHERE id = 2;


SELECT c.courses_id, c.user_id, c.course_name, c.active, 
(SELECT u.name from users u WHERE c.user_id = u.id) 
from user_courses c 
WHERE courses_id = 7;



SELECT c.courses_id, c.user_id, c.course_name, c.active, 
(SELECT u.name from users u WHERE c.user_id = u.id) 
from user_courses c 
WHERE user_id = 10;



