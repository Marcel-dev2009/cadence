import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});
  @override
  State<HomeScreen> createState() => _MyHomeScreen();
}

class _MyHomeScreen extends State<HomeScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0B0B0C),
      body: SafeArea(
        child: Column(
          children: [
            // Header section
            Padding(
              padding: const EdgeInsets.symmetric(horizontal:16.0, vertical: 8.00),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                crossAxisAlignment: CrossAxisAlignment.center,
               children: [
                // App title
                const Text(
                  'Cadence',
                  style: TextStyle(
                    color: Color(0xFFF5F5F5), // whitesmoke
                    fontSize: 22,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Row(
                  mainAxisSize: MainAxisSize.min,
                  spacing: -8.0,
                  children: [
                    IconButton(
                     icon: const Icon(
                      Icons.search,
                      color: Colors.white,
                     ), 
                     onPressed: () {
                      // Function to handle search
                     },
                    ),
                    IconButton(
                      icon: const Icon(Icons.notifications_none_rounded, color: Colors.white,),
                      onPressed: () {
                        // Function to notification screen
                      },
                    ),
                    IconButton(
                      icon: const Icon(Icons.settings, color:Colors.white),
                      onPressed: () {
                        // route to settings screen
                      },
                    )
                  ],
                )
               ],
              ),
            )
          ],
        ),
      )
    );
  }
}
