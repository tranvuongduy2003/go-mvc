# Documentation Migration Guide

## 📖 What Changed?

The Go MVC documentation has been **reorganized into a professional ebook format** with improved structure, navigation, and content quality.

### Old Structure
```
docs/
├── AI_API_GENERATION_RULES.md
├── AI_QUICK_START.md
├── API.md
├── ARCHITECTURE.md
├── DEVELOPMENT.md
├── MIGRATIONS.md
├── ...and 13 more files
```

### New Structure
```
docs/
├── BOOK.md                    # 📚 Start here!
├── README.md                  # Overview & guide
├── INDEX.md                   # Quick access
├── book/                      # All content organized
│   ├── 01-getting-started/
│   ├── 02-architecture/
│   ├── 03-development-guide/
│   ├── 04-features/
│   ├── 05-ai-development/
│   └── 06-operations/
└── [original files preserved]
```

## 🚀 Quick Start with New Documentation

### Step 1: Start with the Book
```bash
# Open the main book
open docs/BOOK.md
```

The book contains:
- Complete table of contents
- Reading paths for different roles
- Navigation guide
- All chapters organized

### Step 2: Choose Your Path

**I'm new to Go MVC**
→ Read [Part I: Getting Started](book/01-getting-started/)

**I want to understand architecture**
→ Jump to [Part II: Architecture](book/02-architecture/)

**I'm ready to code**
→ Check [Part III: Development Guide](book/03-development-guide/)

**I need a specific feature**
→ Browse [Part IV: Features](book/04-features/)

**I want AI code generation**
→ See [Part V: AI Development](book/05-ai-development/)

**I'm deploying to production**
→ Read [Part VI: Operations](book/06-operations/)

### Step 3: Use Quick Access

For fast navigation, use:
- **INDEX.md**: Topic-based quick links
- **README.md**: Overview and status
- **Book chapters**: Directly open relevant chapter

## 🔍 Finding What You Need

### Old File → New Location

| Old File | New Location | Notes |
|----------|-------------|-------|
| `AI_QUICK_START.md` | `book/05-ai-development/01-ai-quick-start.md` | English, organized |
| `AI_API_GENERATION_RULES.md` | `book/05-ai-development/02-api-generation-rules.md` | English, updated |
| `CODE_GENERATION_GUIDELINES.md` | `book/05-ai-development/03-code-generation-guidelines.md` | English, updated |
| `ARCHITECTURE.md` | `book/02-architecture/01-architecture-overview.md` | Organized |
| `PROJECT_STRUCTURE.md` | `book/02-architecture/02-project-structure.md` | Organized |
| `DEPENDENCY_INJECTION.md` | `book/02-architecture/07-dependency-injection.md` | Organized |
| `ARCHITECTURE_EXAMPLES.md` | `book/02-architecture/08-architecture-examples.md` | Organized |
| `DEVELOPMENT.md` | `book/03-development-guide/01-development-workflow.md` | Organized |
| `MIGRATIONS.md` | `book/03-development-guide/02-migrations.md` | Organized |
| `API.md` | `book/03-development-guide/04-api-development.md` | Organized |
| `RBAC_USAGE.md` | `book/04-features/01-authentication.md` | Organized |
| `BACKGROUND_JOBS.md` | `book/04-features/02-background-jobs.md` | Organized |
| `EMAIL_SERVICE.md` | `book/04-features/03-email-service.md` | Organized |
| `FILE_UPLOAD.md` | `book/04-features/04-file-upload.md` | Organized |
| `MESSAGE_DEDUPLICATION.md` | `book/04-features/05-message-deduplication.md` | Organized |
| `NATS_MESSAGING.md` | `book/04-features/06-nats-messaging.md` | Organized |
| `TRACING.md` | `book/04-features/07-tracing.md` | Organized |
| `DEPLOYMENT.md` | `book/06-operations/01-deployment.md` | Organized |
| `QUICK_REFERENCE.md` | `book/01-getting-started/03-quick-reference.md` | English, enhanced |

### New Content (Not in Old Docs)

| New Chapter | Location | Description |
|-------------|----------|-------------|
| **Introduction** | `book/01-getting-started/01-introduction.md` | Complete project overview |
| **Quick Start** | `book/01-getting-started/02-quick-start.md` | 15-minute setup guide |
| **Domain Layer** | `book/02-architecture/03-domain-layer.md` | Deep dive into domain layer |

## 💡 Benefits of New Structure

### Better Organization
- ✅ Logical progression from basics to advanced
- ✅ Grouped by topic and complexity
- ✅ Clear navigation paths
- ✅ Standalone chapters

### Improved Content
- ✅ Professional ebook format
- ✅ Comprehensive new chapters
- ✅ Enhanced examples
- ✅ Better formatting

### Easy Navigation
- ✅ Table of contents in BOOK.md
- ✅ Quick index in INDEX.md
- ✅ Cross-chapter references
- ✅ Multiple reading paths

### English Content
- ✅ All new content in English
- ✅ Updated AI generation docs to English
- ✅ Consistent terminology
- ✅ Professional language

## 🔄 Updating Your Bookmarks

### If You Bookmarked Old Files

**Old bookmark**: `docs/AI_QUICK_START.md`
**New bookmark**: `docs/book/05-ai-development/01-ai-quick-start.md`

**Old bookmark**: `docs/ARCHITECTURE.md`
**New bookmark**: `docs/BOOK.md` (for TOC) or `docs/book/02-architecture/01-architecture-overview.md`

**Old bookmark**: `docs/DEVELOPMENT.md`
**New bookmark**: `docs/book/03-development-guide/01-development-workflow.md`

### Recommended Bookmarks

Save these for quick access:
- `docs/BOOK.md` - Main entry point
- `docs/INDEX.md` - Quick navigation
- `docs/book/01-getting-started/02-quick-start.md` - Quick reference
- `docs/book/05-ai-development/01-ai-quick-start.md` - AI features

## 🤝 For Contributors

### Adding New Content

1. **Determine the right part**
   - Getting Started: Basics and onboarding
   - Architecture: Design and structure
   - Development Guide: Development workflows
   - Features: Specific functionality
   - AI Development: Code generation
   - Operations: Deployment and ops

2. **Create the chapter file**
   ```bash
   # Use naming convention: NN-chapter-name.md
   touch docs/book/02-architecture/04-application-layer.md
   ```

3. **Update BOOK.md**
   - Add chapter to table of contents
   - Update part description if needed

4. **Update INDEX.md**
   - Add quick link if it's a common topic

5. **Follow formatting standards**
   - Use clear headings
   - Include code examples
   - Add cross-references
   - Use consistent style

### Translating Content

If you find Vietnamese content:

1. **Translate to English**
   - Use clear, professional language
   - Maintain technical accuracy
   - Keep code examples as-is

2. **Update cross-references**
   - Ensure links point to new locations
   - Update relative paths

3. **Test links**
   - Verify all cross-references work
   - Check code examples

## 📞 Getting Help

### Documentation Issues

**Found a broken link?**
- Check the [mapping table](#old-file--new-location) above
- Open an issue if the file is missing

**Can't find a topic?**
- Check INDEX.md for quick links
- Search in BOOK.md table of contents
- Open an issue to request clarification

**Content in Vietnamese?**
- Some organized files may still have Vietnamese content
- Original files (in root docs/) are preserved for reference
- Translations in progress

### Questions?

- 📖 Read the [README.md](README.md) for overview
- 🔍 Check [INDEX.md](INDEX.md) for quick access
- 📚 Browse [BOOK.md](BOOK.md) for full content
- 🐛 [Open an issue](https://github.com/tranvuongduy2003/go-mvc/issues) for help

## ✅ Checklist for Users

- [ ] Bookmark `docs/BOOK.md` as main entry point
- [ ] Update any saved links to old documentation
- [ ] Explore the new structure
- [ ] Find your relevant reading path
- [ ] Report any issues or suggestions

---

## 🎉 Enjoy the New Documentation!

The new ebook format makes it easier to:
- Learn Go MVC systematically
- Find information quickly
- Understand complex topics
- Reference while coding
- Onboard new team members

**Start reading**: [BOOK.md](BOOK.md)
