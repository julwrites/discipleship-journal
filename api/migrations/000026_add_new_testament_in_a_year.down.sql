-- Remove New Testament in a Year Reading Plan
DELETE FROM reading_plan_days WHERE reading_plan_id = '7d659c37-4f3b-47b4-ba41-a7a90d68d7fb';
DELETE FROM reading_plans WHERE id = '7d659c37-4f3b-47b4-ba41-a7a90d68d7fb';
