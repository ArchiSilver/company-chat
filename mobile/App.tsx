/**
 * Company Chat Mobile App
 * React Native application for company messaging
 *
 * @format
 */

import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import React from 'react';
import { StatusBar } from 'react-native';
import { Provider as PaperProvider } from 'react-native-paper';
import { SafeAreaProvider } from 'react-native-safe-area-context';

// Types
import { RootStackParamList } from './src/types/navigation';

// Screens
import ChatListScreen from './src/screens/ChatListScreen';
import ChatScreen from './src/screens/ChatScreen';
import LoginScreen from './src/screens/LoginScreen';
import ProfileScreen from './src/screens/ProfileScreen';

// Context
import { AuthProvider } from './src/context/AuthContext';
import { ChatProvider } from './src/context/ChatContext';

const Stack = createStackNavigator<RootStackParamList>();

function App() {
  return (
    <SafeAreaProvider>
      <PaperProvider>
        <AuthProvider>
          <ChatProvider>
            <NavigationContainer>
              <StatusBar barStyle="light-content" backgroundColor="#007AFF" />
              <Stack.Navigator
                initialRouteName="Login"
                screenOptions={{
                  headerStyle: {
                    backgroundColor: '#007AFF',
                  },
                  headerTintColor: '#fff',
                  headerTitleStyle: {
                    fontWeight: 'bold',
                  },
                }}>
                <Stack.Screen
                  name="Login"
                  component={LoginScreen}
                  options={{ headerShown: false }}
                />
                <Stack.Screen
                  name="ChatList"
                  component={ChatListScreen}
                  options={{ title: 'Chats' }}
                />
                <Stack.Screen
                  name="Chat"
                  component={ChatScreen}
                  options={({ route }) => ({
                    title: route.params?.chatName || 'Chat',
                  })}
                />
                <Stack.Screen
                  name="Profile"
                  component={ProfileScreen}
                  options={{ title: 'Profile' }}
                />
              </Stack.Navigator>
            </NavigationContainer>
          </ChatProvider>
        </AuthProvider>
      </PaperProvider>
    </SafeAreaProvider>
  );
}

export default App;
