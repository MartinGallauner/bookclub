import 'dart:developer';

import 'package:bookclub_api/bookclub_api.dart';
import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:frontend/firebase_options.dart';
import 'package:frontend/pages/contacts_page.dart';
import 'package:frontend/pages/library_page.dart';
import 'package:frontend/pages/login_page.dart';
import 'package:frontend/pages/search_page.dart';
import 'package:frontend/services/api_client.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/services/contact_service.dart';
import 'package:frontend/services/library_service.dart';
import 'package:frontend/widgets/add_book_bottom_sheet.dart';
import 'package:frontend/widgets/add_contact_bottom_sheet.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:provider/provider.dart';

import 'layout/app_navigation_rail.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
  runApp(const App());
}

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        Provider(create: (_) => ApiClient()),
        Provider(create: (_) => GoogleSignIn.instance),
        Provider(
          create: (context) => AuthService(
            context.read<ApiClient>(),
            context.read<GoogleSignIn>(),
          ),
        ),
        ChangeNotifierProvider(create: (_) => AppState()),
        ChangeNotifierProxyProvider<AppState, LibraryService?>(
          update: (context, appState, previous) {
            if (appState.user != null) {
              return LibraryService(
                appState.user!.id,
                FirebaseFirestore.instance,
              );
            } else {
              return null;
            }
          },
          create: (BuildContext context) {
            return null;
          },
        ),
        ChangeNotifierProxyProvider<AppState, ContactService?>(
          create: (BuildContext context) {
            return null;
          },
          update: (context, appState, previous) {
            if (appState.user != null) {
              return ContactService(
                appState.user!.id,
                FirebaseFirestore.instance,
              );
            } else {
              return null;
            }
          },
        ),
      ],
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
  UserProfile? user;

  void login(UserProfile user) {
    this.user = user;
    notifyListeners();
  }

  void logout() {
    user = null;
    notifyListeners();
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
        page = ContactsPage();
        break;
      default:
        throw UnimplementedError('no widget for $selectedIndex');
    }

    var user = context.watch<AppState>().user;

    return LayoutBuilder(
      builder: (context, constraints) {
        return Scaffold(
          appBar: AppBar(
            title: Text('BookClub'),
            backgroundColor: Theme.of(context).colorScheme.inversePrimary,
            actions: [
              GestureDetector(
                onTap: () {
                  Clipboard.setData(ClipboardData(text: user?.id ?? ''));
                  ScaffoldMessenger.of(
                    context,
                  ).showSnackBar(SnackBar(content: Text('User ID copied!')));
                },
                child: Padding(
                  padding: EdgeInsets.symmetric(horizontal: 8.0),
                  child: Center(child: Text(user?.id ?? '')),
                ),
              ),
              IconButton(
                icon: Icon(Icons.logout),
                onPressed: () {
                  context.read<AuthService>().signOut();
                  context.read<AppState>().logout();
                },
              ),
            ],
          ),
          body: Row(
            children: [
              AppNavigationRail(
                isExtended: constraints.maxWidth >= 600,
                selectedIndex: selectedIndex,
                onDestinationSelected: (newIndex) {
                  setState(() {
                    selectedIndex = newIndex;
                  });
                },
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
                        isScrollControlled: true,
                        useSafeArea: true,
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
