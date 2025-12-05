# Manual Test Checklist: P1 User Stories

**Feature**: Multilingual Dictionary with Vocabulary Game  
**Date**: 2025-01-27  
**Scope**: Priority 1 (P1) User Stories - MVP Features

## Test Environment Setup

- [x] Backend API running on `http://localhost:9090`
- [x] Frontend running on `http://localhost:5173`
- [x] Database connected and migrations applied
- [x] Sample data loaded (languages, words, topics, levels)

---

## User Story 1: Landing Page Navigation (P1)

**Goal**: Users can navigate from landing page to game list or dictionary lookup.

### Test Steps

1. [x] Navigate to `http://localhost:5173/`
2. [x] Verify landing page displays with title "English Coach"
3. [x] Verify "Play Game" button is visible and clickable
4. [x] Verify "Dictionary Lookup" button is visible and clickable
5. [x] Click "Play Game" button
6. [x] Verify navigation to `/games` page
7. [x] Navigate back to landing page
8. [x] Click "Dictionary Lookup" button
9. [x] Verify navigation to `/dictionary` page

### Expected Results

- ✅ Landing page loads within 2 seconds (SC-001)
- ✅ Both buttons are clearly visible
- ✅ Navigation works correctly for both buttons
- ✅ No console errors

### Issues Found

- [ ] Issue 1: ________________
- [ ] Issue 2: ________________

---

## User Story 2: Game List Display (P1)

**Goal**: Users can see available games and select vocabulary game.

### Test Steps

1. [x] Navigate to `/games` page
2. [x] Verify "Về trang chủ" (Back to Home) button is visible
3. [x] Verify game list displays with "Học Từ Vựng" (Vocabulary Game) option
4. [x] Verify game card shows description
5. [x] Click on vocabulary game card
6. [x] Verify navigation to `/games/vocab/config` page
7. [x] Click "Về trang chủ" button
8. [x] Verify navigation back to landing page

### Expected Results

- ✅ Game list displays correctly
- ✅ Vocabulary game option is visible
- ✅ Clicking game navigates to configuration page
- ✅ Back button works correctly

### Issues Found

- [ ] Issue 1: ________________
- [ ] Issue 2: ________________

---

## User Story 3: Vocabulary Game Configuration (P1)

**Goal**: Users can configure game session and start playing.

### Test Steps - Valid Configuration

1. [x] Navigate to `/games/vocab/config`
2. [x] Verify "Quay Lại" (Back) button is visible
3. [x] Verify language dropdowns load (source and target)
4. [x] Select source language (e.g., English)
5. [x] Select target language (e.g., Vietnamese) - different from source
6. [x] Select mode: "Topic" or "Level"
7. [x] If Topic mode: Select a topic
8. [x] If Level mode: Select a level
9. [x] Click "Bắt Đầu Chơi" (Start Game) button
10. [ ] Verify navigation to game play page with session ID

### Test Steps - Invalid Configuration

1. [ ] Try to select same language for source and target
2. [ ] Verify error message: "Ngôn ngữ nguồn và ngôn ngữ đích phải khác nhau"
3. [ ] Try to start without selecting mode
4. [ ] Verify error message appears
5. [ ] Try to start without selecting topic/level
6. [ ] Verify error message appears

### Expected Results

- ✅ Configuration can be completed in under 30 seconds (SC-002)
- ✅ Validation errors display in Vietnamese (FR-025)
- ✅ Game session is created successfully
- ✅ Navigation to game play page works

### Issues Found

- [ ] Issue 1: ________________
- [ ] Issue 2: ________________

---

## User Story 4: Playing Vocabulary Game (P1)

**Goal**: Users can play game by answering multiple-choice questions.

### Test Steps

1. [ ] Start game session from configuration page
2. [ ] Verify game play page loads
3. [ ] Verify first question displays within 1 second (SC-003)
4. [ ] Verify source word is displayed prominently
5. [ ] Verify 4 answer options (A, B, C, D) are displayed
6. [ ] Click on an answer option (e.g., "A")
7. [ ] Verify answer is submitted
8. [ ] Verify progress indicator updates
9. [ ] Verify next question loads automatically after short delay
10. [ ] Answer all questions in the session
11. [ ] Verify completion screen appears
12. [ ] Verify "View Statistics" button is visible

### Test Steps - Error Handling

1. [ ] Disconnect network during game play
2. [ ] Try to submit an answer
3. [ ] Verify error message: "Không thể kết nối đến máy chủ..."
4. [ ] Reconnect network
5. [ ] Verify answer can be submitted successfully

### Expected Results

- ✅ Questions load within 1 second (SC-003)
- ✅ All questions can be answered
- ✅ Answers are recorded correctly
- ✅ Completion screen appears after all questions
- ✅ Network errors are handled gracefully
- ✅ Error messages are in Vietnamese

### Issues Found

- [ ] Issue 1: ________________
- [ ] Issue 2: ________________

---

## User Story 5: Viewing Game Statistics (P1)

**Goal**: Users can view statistics after completing a game.

### Test Steps

1. [ ] Complete a game session (from User Story 4)
2. [ ] Click "View Statistics" button on completion screen
3. [ ] Verify navigation to statistics page
4. [ ] Verify statistics display:
   - Total questions
   - Correct answers
   - Wrong answers
   - Accuracy percentage
   - Session duration
   - Average response time
5. [ ] Verify "Chơi lại" (Play Again) button is visible
6. [ ] Verify "Quay lại danh sách game" (Back to Game List) button is visible
7. [ ] Click "Chơi lại" button
8. [ ] Verify navigation to configuration page
9. [ ] Navigate back to statistics page
10. [ ] Click "Quay lại danh sách game" button
11. [ ] Verify navigation to game list page

### Expected Results

- ✅ Statistics display correctly
- ✅ All statistics fields are accurate
- ✅ Navigation buttons work correctly
- ✅ Statistics page loads without errors

### Issues Found

- [ ] Issue 1: ________________
- [ ] Issue 2: ________________

---

## User Story 6: Dictionary Lookup (P2)

**Goal**: Users can search for words and view detailed information.

### Test Steps - Search

1. [ ] Navigate to `/dictionary` page
2. [ ] Verify "Về trang chủ" (Back to Home) button is visible
3. [ ] Verify search input field is visible
4. [ ] Enter a search query (e.g., "hello")
5. [ ] Verify search results appear after debounce delay
6. [ ] Verify results show word lemma and romanization
7. [ ] Verify pagination controls appear if results > 20
8. [ ] Click on a word in results
9. [ ] Verify navigation to word detail page

### Test Steps - Word Detail

1. [ ] Navigate to word detail page
2. [ ] Verify "Quay lại tìm kiếm" (Back to Search) button is visible
3. [ ] Verify word information displays:
   - Word lemma
   - Romanization (if available)
   - Pronunciations (if available)
   - Senses with definitions
   - Translations for each sense
   - Examples (if available)
4. [ ] Click "Quay lại tìm kiếm" button
5. [ ] Verify navigation back to dictionary search page

### Test Steps - Empty Results

1. [ ] Search for a word that doesn't exist (e.g., "xyzabc123")
2. [ ] Verify message: "Không tìm thấy từ nào phù hợp với..."
3. [ ] Verify no errors occur

### Expected Results

- ✅ Dictionary lookup returns results within 1 second for 95% of searches (SC-005)
- ✅ Search results are relevant
- ✅ Word detail page shows complete information
- ✅ Empty results are handled gracefully
- ✅ Pagination works correctly

### Issues Found

- [ ] Issue 1: ________________
- [ ] Issue 2: ________________

---

## Cross-Cutting Concerns

### Error Handling

- [ ] All error messages are in Vietnamese (FR-025)
- [ ] Network errors display user-friendly messages
- [ ] Validation errors are clear and actionable
- [ ] 404 errors handled gracefully

### Navigation

- [ ] All pages have back/home navigation buttons
- [ ] Navigation is consistent across pages
- [ ] Browser back button works correctly

### Performance

- [ ] Landing page loads within 2 seconds (SC-001)
- [ ] Game configuration completes in under 30 seconds (SC-002)
- [ ] Questions load within 1 second (SC-003)
- [ ] Dictionary search returns results within 1 second (SC-005)

### Input Validation

- [ ] All form inputs are validated
- [ ] Validation errors display immediately
- [ ] Invalid data cannot be submitted

---

## Overall Test Summary

**Test Date**: _______________  
**Tester**: _______________  
**Environment**: Development / Staging / Production

### Test Results

- **Total Test Cases**: ___
- **Passed**: ___
- **Failed**: ___
- **Blocked**: ___

### Critical Issues

1. ________________
2. ________________
3. ________________

### Recommendations

1. ________________
2. ________________
3. ________________

### Sign-off

- [ ] All P1 user stories are functional
- [ ] Critical issues resolved
- [ ] Ready for next phase / deployment

**Approved by**: _______________  
**Date**: _______________

