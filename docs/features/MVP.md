# MVP Feature Specifications

## 1. Authentication
- **Sign-up/Log-in**: Users must authenticate to access the app.
- **Provider**: Google Sign-in via Firebase Auth.
- **Session**: Persisted session (standard PWA behavior).

## 2. Dashboard
- **Layout**: Main landing view after login.
- **Content**:
  - List of Journal Notes displayed as cards.
  - Ordered by Most Recently Used (MRU).
- **Interactions**:
  - **Search Bar**: Sticky at the bottom. Filters notes by title/content.
  - **FAB (+)**: Bottom right button to create a new Journal Note.
  - **Chat Entry**: Button to start a new AI Chat session.

## 3. User Settings
- **Location**: Accessible from Dashboard (e.g., profile icon).
- **Fields**:
  - **Username**: Globally unique identifier.
  - **Preferred Bible Version**: Selection (e.g., NIV, ESV, KJV).

## 4. Journaling
- **Core Entity**: "Journal Note".
- **Fields**: Title, Content, Date, Bible Context.
- **Editor**:
  - Markdown support.
  - Rich media support: Images, GIFs, Emojis.
- **Context**:
  - User can select Bible passages to associate with the note.
- **AI Assistance**:
  - "Ask AI": User can ask a question about the note's context (Bible passages + content).
  - *Constraint*: AI ignores images/GIFs in the context.

## 5. AI Chat Session
- **Workflow**:
  1.  **Context Selection**: User selects Bible passages.
  2.  **Themes (Optional)**: User enters keywords/themes.
  3.  **Prompt**: User enters the initial question/prompt.
  4.  **Processing**: System sends context + prompt to BibleAIAPI.
  5.  **Output**: A new Journal Note is created containing:
      - The selected passages.
      - The prompt.
      - The AI's response.
      - Title generated from the prompt.
