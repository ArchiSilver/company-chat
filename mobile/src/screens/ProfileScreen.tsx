import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import React from 'react';
import { Alert, StyleSheet, View } from 'react-native';
import { Avatar, Button, Card, Text, Title } from 'react-native-paper';
import { useAuth } from '../context/AuthContext';
import { RootStackParamList } from '../types/navigation';

type ProfileScreenNavigationProp = StackNavigationProp<RootStackParamList, 'Profile'>;

const ProfileScreen: React.FC = () => {
    const { user, logout } = useAuth();
    const navigation = useNavigation<ProfileScreenNavigationProp>();

    const handleLogout = () => {
        Alert.alert(
            'Logout',
            'Are you sure you want to logout?',
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Logout',
                    style: 'destructive',
                    onPress: async () => {
                        await logout();
                        navigation.replace('Login');
                    },
                },
            ],
        );
    };

    if (!user) {
        return (
            <View style={styles.centerContainer}>
                <Text>User not found</Text>
            </View>
        );
    }

    return (
        <View style={styles.container}>
            <Card style={styles.profileCard}>
                <Card.Content style={styles.profileContent}>
                    <Avatar.Text
                        size={80}
                        label={user.username.charAt(0).toUpperCase()}
                        style={styles.avatar}
                    />

                    <Title style={styles.username}>{user.username}</Title>
                    <Text style={styles.email}>{user.email}</Text>
                    <Text style={styles.role}>Role: {user.role}</Text>

                    <Text style={styles.joinDate}>
                        Joined: {new Date(user.created_at).toLocaleDateString()}
                    </Text>
                </Card.Content>
            </Card>

            <View style={styles.actionsContainer}>
                <Button
                    mode="outlined"
                    onPress={() => navigation.goBack()}
                    style={styles.button}>
                    Back to Chats
                </Button>

                <Button
                    mode="contained"
                    onPress={handleLogout}
                    style={[styles.button, styles.logoutButton]}
                    buttonColor="#FF3B30">
                    Logout
                </Button>
            </View>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#f5f5f5',
        padding: 16,
    },
    centerContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    profileCard: {
        elevation: 4,
        marginBottom: 24,
    },
    profileContent: {
        alignItems: 'center',
        padding: 24,
    },
    avatar: {
        marginBottom: 16,
    },
    username: {
        fontSize: 24,
        marginBottom: 8,
    },
    email: {
        fontSize: 16,
        color: '#666',
        marginBottom: 8,
    },
    role: {
        fontSize: 14,
        color: '#666',
        marginBottom: 16,
    },
    joinDate: {
        fontSize: 12,
        color: '#999',
    },
    actionsContainer: {
        gap: 12,
    },
    button: {
        marginBottom: 8,
    },
    logoutButton: {
        marginTop: 16,
    },
});

export default ProfileScreen;