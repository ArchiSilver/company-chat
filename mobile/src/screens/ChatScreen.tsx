import { RouteProp, useRoute } from '@react-navigation/native';
import React, { useEffect, useRef, useState } from 'react';
import {
    FlatList,
    KeyboardAvoidingView,
    Platform,
    StyleSheet,
    View,
} from 'react-native';
import { ActivityIndicator, Avatar, Card, Text, TextInput } from 'react-native-paper';
import { useAuth } from '../context/AuthContext';
import { useChat } from '../context/ChatContext';
import { Message } from '../types';
import { RootStackParamList } from '../types/navigation';

type ChatScreenRouteProp = RouteProp<RootStackParamList, 'Chat'>;

const ChatScreen: React.FC = () => {
    const [messageText, setMessageText] = useState('');
    const [isSending, setIsSending] = useState(false);
    const { currentChat, messages, isLoading, sendMessage, connectWebSocket } = useChat();
    const { user } = useAuth();
    const route = useRoute<ChatScreenRouteProp>();
    const flatListRef = useRef<FlatList>(null);

    const { chatId } = route.params;

    useEffect(() => {
        // Connect to WebSocket for real-time updates
        connectWebSocket(chatId);
    }, [chatId, connectWebSocket]);

    const handleSendMessage = async () => {
        if (!messageText.trim() || isSending) return;

        setIsSending(true);
        try {
            await sendMessage(messageText.trim());
            setMessageText('');
            // Scroll to bottom after sending
            setTimeout(() => {
                flatListRef.current?.scrollToEnd({ animated: true });
            }, 100);
        } catch (error) {
            console.error('Error sending message:', error);
        } finally {
            setIsSending(false);
        }
    };

    const renderMessage = ({ item }: { item: Message }) => {
        const isOwnMessage = item.user_id === user?.id;

        return (
            <View style={[
                styles.messageContainer,
                isOwnMessage ? styles.ownMessage : styles.otherMessage
            ]}>
                {!isOwnMessage && (
                    <Avatar.Text
                        size={32}
                        label={item.user?.username?.charAt(0).toUpperCase() || 'U'}
                        style={styles.messageAvatar}
                    />
                )}
                <Card style={[
                    styles.messageCard,
                    isOwnMessage ? styles.ownMessageCard : styles.otherMessageCard
                ]}>
                    <Card.Content style={styles.messageContent}>
                        {!isOwnMessage && (
                            <Text style={styles.messageUser}>{item.user?.username}</Text>
                        )}
                        <Text style={styles.messageText}>{item.content}</Text>
                        <Text style={styles.messageTime}>
                            {new Date(item.created_at).toLocaleTimeString([], {
                                hour: '2-digit',
                                minute: '2-digit'
                            })}
                        </Text>
                    </Card.Content>
                </Card>
            </View>
        );
    };

    const renderFooter = () => {
        if (isLoading) {
            return (
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="small" />
                </View>
            );
        }
        return null;
    };

    if (!currentChat) {
        return (
            <View style={styles.centerContainer}>
                <ActivityIndicator size="large" />
            </View>
        );
    }

    return (
        <KeyboardAvoidingView
            style={styles.container}
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
            keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}>
            <FlatList
                ref={flatListRef}
                data={messages}
                keyExtractor={item => item.id}
                renderItem={renderMessage}
                style={styles.messagesList}
                contentContainerStyle={styles.messagesContainer}
                onContentSizeChange={() => flatListRef.current?.scrollToEnd({ animated: false })}
                onLayout={() => flatListRef.current?.scrollToEnd({ animated: false })}
                ListFooterComponent={renderFooter}
            />

            <View style={styles.inputContainer}>
                <TextInput
                    value={messageText}
                    onChangeText={setMessageText}
                    placeholder="Type a message..."
                    mode="outlined"
                    multiline
                    style={styles.messageInput}
                    right={
                        <TextInput.Icon
                            icon="send"
                            onPress={handleSendMessage}
                            disabled={!messageText.trim() || isSending}
                        />
                    }
                />
            </View>
        </KeyboardAvoidingView>
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
    },
    messagesList: {
        flex: 1,
    },
    messagesContainer: {
        padding: 16,
    },
    messageContainer: {
        flexDirection: 'row',
        marginBottom: 12,
        alignItems: 'flex-end',
    },
    ownMessage: {
        justifyContent: 'flex-end',
    },
    otherMessage: {
        justifyContent: 'flex-start',
    },
    messageAvatar: {
        marginRight: 8,
    },
    messageCard: {
        maxWidth: '70%',
        elevation: 1,
    },
    ownMessageCard: {
        backgroundColor: '#007AFF',
    },
    otherMessageCard: {
        backgroundColor: '#fff',
    },
    messageContent: {
        padding: 8,
    },
    messageUser: {
        fontSize: 12,
        fontWeight: 'bold',
        marginBottom: 4,
        color: '#666',
    },
    messageText: {
        fontSize: 16,
        lineHeight: 20,
    },
    messageTime: {
        fontSize: 10,
        color: '#999',
        alignSelf: 'flex-end',
        marginTop: 4,
    },
    inputContainer: {
        padding: 16,
        backgroundColor: '#fff',
        borderTopWidth: 1,
        borderTopColor: '#e0e0e0',
    },
    messageInput: {
        backgroundColor: '#fff',
    },
    loadingContainer: {
        padding: 16,
        alignItems: 'center',
    },
});

export default ChatScreen;