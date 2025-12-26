import 'package:flutter/material.dart';
import 'package:frontend/pages/library_page.dart';
import 'package:frontend/pages/network_page.dart';
import 'package:frontend/pages/search_page.dart';
import 'package:frontend/widgets/add_book_bottom_sheet.dart';
import 'package:frontend/widgets/add_contact_bottom_sheet.dart';
import 'package:provider/provider.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => MyAppState(),
      child: MaterialApp(
        title: 'bookclub',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigoAccent),
        ),
        home: MyHomePage(),
      ),


    );
  }
}

class MyAppState extends ChangeNotifier {

}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key});

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage>{
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
              title: Text('Book Club'),
              backgroundColor: Theme.of(context).colorScheme.inversePrimary,
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
                    print('Add book pressed');
                    showModalBottomSheet(
                        context: context,
                        builder: (context) => AddBookBottomSheet(),
                    );
                  } else if (selectedIndex == 2) {
                    print('Add contact pressed');
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
            }
          );
        }
  }










