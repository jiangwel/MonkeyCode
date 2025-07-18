ALTER TABLE tasks ADD column is_suggested boolean default false;
ALTER TABLE task_records ADD column code_lines int default 0;
ALTER TABLE task_records ADD column code text;