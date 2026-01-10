import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import React, { useEffect } from 'react';
import {
    Alert,
    FlatList,
    StyleSheet,
    TouchableOpacity,
    View,
} from 'react-native';
import { ActivityIndicator, Avatar, Card, FAB, Text, Title } from 'react-native-paper';
import { useAuth } from '../context/AuthContext';
import { useChat } from '../context/ChatContext';
import { Chat } from '../types';
import { RootStackParamList } from '../types/navigation';

type ChatListScreenNavigationProp = StackNavigationProp<RootStackParamList, 'ChatList'>;

const ChatListScreen: React.FC = () => {
    const { chats, isLoading, error, loadChats } = useChat();
    const { user, logout } = useAuth();
    const navigation = useNavigation<ChatListScreenNavigationProp>();

    useEffect(() => {
        loadChats();
    }, [loadChats]);

    const handleChatPress = (chat: Chat) => {
        navigation.navigate('Chat', {
            chatId: chat.id,
            chatName: chat.name,
        });
    };

    const handleProfilePress = () => {
        navigation.navigate('Profile');
    };

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

    const renderChatItem = ({ item }: { item: Chat }) => {
        const lastMessage = item.last_message;
        const chatName = item.name || item.participants
            .filter(p => p.id !== user?.id)
            .map(p => p.username)
            .join(', ');

        return (
            <TouchableOpacity onPress={() => handleChatPress(item)}>
                <Card style={styles.chatCard}>
                    <Card.Content style={styles.chatContent}>
                        <Avatar.Text
                            size={50}
                            label={chatName.charAt(0).toUpperCase()}
                            style={styles.avatar}
                        />
                        <View style={styles.chatInfo}>
                            <Title style={styles.chatTitle}>{chatName}</Title>
                            {lastMessage && (
                                <Text style={styles.lastMessage} numberOfLines={1}>
                                    {lastMessage.user?.username}: {lastMessage.content}
                                </Text>
                            )}
                        </View>
                    </Card.Content>
                </Card>
            </TouchableOpacity>
        );
    };

    if (isLoading && chats.length === 0) {
        return (
            <View style={styles.centerContainer}>
                <ActivityIndicator size="large" />
            </View>
        );
    }

    if (error) {
        return (
            <View style={styles.centerContainer}>
                <Text style={styles.errorText}>{error}</Text>
                <TouchableOpacity onPress={loadChats} style={styles.retryButton}>
                    <Text style={styles.retryText}>Retry</Text>
                </TouchableOpacity>
            </View>
        );
    }

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Title style={styles.headerTitle}>Chats</Title>
                <TouchableOpacity onPress={handleProfilePress}>
                    <Avatar.Text
                        size={40}
                        label={user?.username?.charAt(0).toUpperCase() || 'U'}
                    />
                </TouchableOpacity>
            </View>

            {chats.length === 0 ? (
                <View style={styles.emptyContainer}>
                    <Text style={styles.emptyText}>No chats yet</Text>
                    <Text style={styles.emptySubtext}>Start a conversation to see your chats here</Text>
                </View>
            ) : (
                <FlatList
                    data={chats}
                    keyExtractor={item => item.id}
                    renderItem={renderChatItem}
                    style={styles.chatList}
                />
            )}

            <FAB
                icon="logout"
                style={styles.fab}
                onPress={handleLogout}
            />
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#f5f5f5',
    },
    centerContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#f5f5f5',
    },
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: 16,
        backgroundColor: '#fff',
        elevation: 2,
    },
    headerTitle: {
        fontSize: 24,
    },
    chatList: {
        flex: 1,
    },
    chatCard: {
        margin: 8,
        elevation: 2,
    },
    chatContent: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    avatar: {
        marginRight: 16,
    },
    chatInfo: {
        flex: 1,
    },
    chatTitle: {
        fontSize: 16,
        marginBottom: 4,
    },
    lastMessage: {
        fontSize: 14,
        color: '#666',
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: 32,
    },
    emptyText: {
        fontSize: 18,
        fontWeight: 'bold',
        marginBottom: 8,
    },
    emptySubtext: {
        fontSize: 14,
        color: '#666',
        textAlign: 'center',
    },
    errorText: {
        fontSize: 16,
        color: 'red',
        marginBottom: 16,
        textAlign: 'center',
    },
    retryButton: {
        padding: 12,
        backgroundColor: '#007AFF',
        borderRadius: 8,
    },
    retryText: {
        color: '#fff',
        fontSize: 16,
    },
    fab: {
        position: 'absolute',
        margin: 16,
        right: 0,
        bottom: 0,
        backgroundColor: '#FF3B30',
    },
});

export default ChatListScreen;