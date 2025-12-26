# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Flutter web frontend for a bookclub application. The app uses a responsive navigation rail pattern with Provider for state management.

### Project Goal

Bookclub allows users to scan the barcodes of their physical books to build their personal library digitally. Users can then make their library searchable for their friends, enabling book discovery and sharing within their social network.

## Current Status

### ✅ What's Already Implemented

**UI Foundation**:
- NavigationRail-based navigation with 3 pages
- Responsive layout (extends at 600px+ breakpoint)
- Provider state management setup (MyAppState)
- Theme with ColorScheme.fromSeed (indigoAccent)

**Pages**:
- Library Page: GridView displaying user's book collection with mock data
- Network Page: GridView showing contacts/friends with mock data (12 pirate-themed test contacts)
- Search Page: Search interface for finding books

**Widgets**:
- AddBookBottomSheet: Bottom sheet for adding new books to library
- AddContactBottomSheet: Bottom sheet for adding new contacts with search functionality
- ContactCard: Card widget for displaying contact information
- Book cards: Grid cards for displaying books in library

**Current Data**: All pages use mock/test data (no backend connection yet)

### 🚧 Planned: Firebase Migration

**Backend Transition**:
- **Removing**: The existing Go backend will be removed
- **Replacing with**: Firebase serverless infrastructure

**Firebase Services to Integrate**:
1. **Firebase Auth**: User authentication and identity management
2. **Cloud Firestore**: NoSQL database for users, books, friend connections
3. **Firebase App Check**: Security and abuse prevention

**Migration Benefits**:
- No server infrastructure to maintain
- Offline-first with real-time synchronization
- Automatic scaling
- Built-in security rules

**Next Steps**:
1. Set up Firebase project and FlutterFire packages
2. Implement Firebase Authentication
3. Design Firestore data model (users, books, friendships)
4. Replace mock data with Firestore queries
5. Add App Check for security
6. Remove Go backend dependencies

## Common Commands

### Development
```bash
# Run the app locally on web (Chrome)
flutter run -d web-server --web-port=8080

# Or use mise task
mise start

# Run the app with hot reload
flutter run -d chrome
```

### Build & Deploy
```bash
# Build for web production
flutter build web

# Build with web renderer options
flutter build web --web-renderer canvaskit  # Better performance for complex UIs
flutter build web --web-renderer html       # Smaller download size
```

### Testing
```bash
# Run all tests
flutter test

# Run a specific test file
flutter test test/widget_test.dart

# Run tests with coverage
flutter test --coverage
```

### Code Quality
```bash
# Analyze code for issues
flutter analyze

# Format code
flutter format .

# Fix common issues
dart fix --apply
```

### Dependencies
```bash
# Get dependencies
flutter pub get

# Update dependencies
flutter pub upgrade

# Check for outdated packages
flutter pub outdated
```

## Architecture

### Navigation Pattern
The app uses a **NavigationRail-based architecture** (main.dart:38-96):
- `MyHomePage` is a StatefulWidget that manages navigation state
- `selectedIndex` controls which page is displayed
- NavigationRail is **responsive**: extends when screen width >= 600px (main.dart:66)
- Pages are swapped via a switch statement (main.dart:44-53)
- Layout uses Row with SafeArea for the rail and Expanded for page content

When adding new pages:
1. Add a new NavigationRailDestination to the destinations list (main.dart:67-76)
2. Add a corresponding case in the switch statement (main.dart:44-53)
3. Create the page widget in `lib/pages/`

### State Management
Uses **Provider** pattern (pubspec.yaml:37):
- `MyAppState` extends ChangeNotifier (main.dart:29-31)
- Currently empty but set up as the app's global state container
- Wrapped around MaterialApp in main.dart:14-25

To add global state:
1. Add properties and methods to MyAppState
2. Call `notifyListeners()` when state changes
3. Access with `context.watch<MyAppState>()` or `context.read<MyAppState>()`

### Page Structure
Pages are organized in `lib/pages/`:
- Each page is a standalone widget
- LibraryPage demonstrates GridView with responsive maxCrossAxisExtent (library_page.dart:23-28)
- Currently using mock data for development (library_page.dart:10-17)

### Responsive Design Patterns
- NavigationRail extends at 600px+ breakpoint
- GridView uses maxCrossAxisExtent (200px) for fluid grid layout
- LayoutBuilder used for constraint-based responsive behavior (main.dart:55)

## Project Structure

```
lib/
├── main.dart              # App entry, navigation setup, global state
└── pages/                 # Page widgets
    └── library_page.dart  # Library grid view with books
```

## Development Notes

### Adding New Navigation Destinations
The navigation uses an index-based system. When adding destinations, ensure the switch statement cases match the destination order to prevent UnimplementedError (main.dart:52).

### Theme & Styling
- Theme uses ColorScheme.fromSeed with indigoAccent (main.dart:19)
- Access theme colors via `Theme.of(context).colorScheme`
- AppBar uses inversePrimary background (main.dart:60)
- Primary container color used for page background (main.dart:87)