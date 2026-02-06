import 'dart:async';
import 'dart:developer';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:frontend/pages/library_page.dart';
import 'package:frontend/pages/login_page.dart';
import 'package:frontend/pages/network_page.dart';
import 'package:frontend/pages/search_page.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/services/library_service.dart';
import 'package:frontend/widgets/add_book_bottom_sheet.dart';
import 'package:frontend/widgets/add_contact_bottom_sheet.dart';
import 'package:provider/provider.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:frontend/firebase_options.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
  runApp(const App());
}

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => AppState(),
      builder: (context, child) {
        return MaterialApp(
          title: 'bookclub',
          debugShowCheckedModeBanner: false,
          theme: ThemeData(
            colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigoAccent),
          ),
          //if user is not logged in, send to them to the Login Page.
          home: context.watch<AppState>().user == null
              ? LoginPage()
              : HomePage(),
        );
      },
    );
  }
}

class AppState extends ChangeNotifier {
  User? user;
  late final StreamSubscription<User?> _authSubscription;
  final LibraryService libraryService = LibraryService();

  AppState() {
    _authSubscription = AuthService().authStateChanges().listen((newUser) {
      user = newUser;
      notifyListeners(); //causes all widgets watching this state to rebuild when the auth state changes
    });
  }

  @override
  void dispose() {
    //we need the dispose to clean up the connection we made above when calling .listen()
    _authSubscription.cancel();
    super.dispose();
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  var selectedIndex = 0;

  @override
  Widget build(BuildContext context) {
    Widget page;
    switch (selectedIndex) {
      case 0:
        page = LibraryPage();
        break;
      case 1:
        page = SearchPage();
        break;
      case 2:
        page = NetworkPage();
        break;
      default:
        throw UnimplementedError('no widget for $selectedIndex');
    }

    return LayoutBuilder(
      builder: (context, constraints) {
        return Scaffold(
          appBar: AppBar(
            title: Text('BookClub'),
            backgroundColor: Theme.of(context).colorScheme.inversePrimary,
            actions: [
              IconButton(onPressed: AuthService().signOut, icon: Icon(Icons.logout))
            ],
          ),
          body: Row(
            children: [
              SafeArea(
                child: NavigationRail(
                  extended: constraints.maxWidth >= 600,
                  destinations: [
                    NavigationRailDestination(
                      icon: Icon(Icons.book),
                      label: Text('Your Library'),
                    ),
                    NavigationRailDestination(
                      icon: Icon(Icons.search),
                      label: Text('Search'),
                    ),
                    NavigationRailDestination(
                      icon: Icon(Icons.person),
                      label: Text('Network'),
                    ),
                  ],
                  selectedIndex: selectedIndex,
                  onDestinationSelected: (value) {
                    setState(() {
                      selectedIndex = value;
                    });
                  },
                ),
              ),
              Expanded(
                child: Container(
                  color: Theme.of(context).colorScheme.primaryContainer,
                  child: page,
                ),
              ),
            ],
          ),
          floatingActionButton: (selectedIndex == 0 || selectedIndex == 2)
              ? FloatingActionButton(
                  onPressed: () {
                    if (selectedIndex == 0) {
                      log('pressed add-book-button');
                      showModalBottomSheet(
                        context: context,
                        builder: (context) => AddBookBottomSheet(),
                      );
                    } else if (selectedIndex == 2) {
                      log('pressed add-contact-button');
                      showModalBottomSheet(
                        context: context,
                        builder: (context) => AddContactBottomSheet(),
                      );
                    }
                  },
                  child: Icon(Icons.add),
                )
              : null,
        );
      },
    );
  }
}
