import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import React, { useState } from 'react';
import {
    Alert,
    KeyboardAvoidingView,
    Platform,
    ScrollView,
    StyleSheet
} from 'react-native';
import { Button, Card, TextInput, Title } from 'react-native-paper';
import { useAuth } from '../context/AuthContext';
import { RootStackParamList } from '../types/navigation';

type LoginScreenNavigationProp = StackNavigationProp<RootStackParamList, 'Login'>;

const LoginScreen: React.FC = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isLogin, setIsLogin] = useState(true);
    const [username, setUsername] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const { login, register } = useAuth();
    const navigation = useNavigation<LoginScreenNavigationProp>();

    const handleSubmit = async () => {
        if (!email || !password) {
            Alert.alert('Error', 'Please fill in all fields');
            return;
        }

        if (!isLogin && !username) {
            Alert.alert('Error', 'Please enter a username');
            return;
        }

        setIsLoading(true);

        try {
            let success = false;

            if (isLogin) {
                success = await login({ email, password });
            } else {
                success = await register({ email, username, password });
            }

            if (success) {
                navigation.replace('ChatList');
            } else {
                Alert.alert('Error', isLogin ? 'Invalid credentials' : 'Registration failed');
            }
        } catch {
            Alert.alert('Error', 'Network error. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    const toggleMode = () => {
        setIsLogin(!isLogin);
        setEmail('');
        setPassword('');
        setUsername('');
    };

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
            <ScrollView contentContainerStyle={styles.scrollContainer}>
                <Card style={styles.card}>
                    <Card.Content>
                        <Title style={styles.title}>
                            {isLogin ? 'Welcome Back' : 'Create Account'}
                        </Title>

                        <TextInput
                            label="Email"
                            value={email}
                            onChangeText={setEmail}
                            mode="outlined"
                            keyboardType="email-address"
                            autoCapitalize="none"
                            autoCorrect={false}
                            style={styles.input}
                        />

                        {!isLogin && (
                            <TextInput
                                label="Username"
                                value={username}
                                onChangeText={setUsername}
                                mode="outlined"
                                autoCapitalize="none"
                                autoCorrect={false}
                                style={styles.input}
                            />
                        )}

                        <TextInput
                            label="Password"
                            value={password}
                            onChangeText={setPassword}
                            mode="outlined"
                            secureTextEntry
                            style={styles.input}
                        />

                        <Button
                            mode="contained"
                            onPress={handleSubmit}
                            loading={isLoading}
                            disabled={isLoading}
                            style={styles.button}>
                            {isLogin ? 'Login' : 'Register'}
                        </Button>

                        <Button
                            mode="text"
                            onPress={toggleMode}
                            style={styles.toggleButton}>
                            {isLogin
                                ? "Don't have an account? Register"
                                : 'Already have an account? Login'}
                        </Button>
                    </Card.Content>
                </Card>
            </ScrollView>
        </KeyboardAvoidingView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#f5f5f5',
    },
    scrollContainer: {
        flexGrow: 1,
        justifyContent: 'center',
        padding: 20,
    },
    card: {
        elevation: 4,
    },
    title: {
        textAlign: 'center',
        marginBottom: 20,
        fontSize: 24,
    },
    input: {
        marginBottom: 16,
    },
    button: {
        marginTop: 16,
        marginBottom: 16,
    },
    toggleButton: {
        marginTop: 8,
    },
});

export default LoginScreen;