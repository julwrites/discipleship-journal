-- Remove Bible in 90 Days Reading Plan
DELETE FROM reading_plan_days WHERE reading_plan_id = 'e93c79f7-f7d9-46d1-ac37-e1d88d83ec02';
DELETE FROM reading_plans WHERE id = 'e93c79f7-f7d9-46d1-ac37-e1d88d83ec02';
