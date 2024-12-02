import React, { useEffect, useRef, useState } from "react";
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  ImageBackground,
  KeyboardAvoidingView,
  Platform,
  Keyboard,
} from "react-native";
import { Feather, Ionicons } from "@expo/vector-icons";
import { Link, router, useLocalSearchParams } from "expo-router";
import { useGlobalContext } from "@/context/GlobalProvider";
import { Chat } from "@/types/Chat";
import axios from "axios";
import { chatServiceHost } from "@/constants/backendUrl";

const ChatScreen = () => {
  const [message, setMessage] = useState("");
  const flatListRef = useRef<FlatList>(null);
  const [keyboardVisible, setKeyboardVisible] = useState(false);

  const { chatId } = useLocalSearchParams();

  const context = useGlobalContext();
  if (context == undefined) throw new Error("Context not defined");

  const { chats, token, currentChat, setCurrentChat } = context;

  useEffect(() => {
    if (!token) {
      router.replace("/");
      return;
    }
    const chat = chats.find(c => c.id.toString() === chatId);
    if (!chat) {
      router.replace("/");
      return;
    }
    setCurrentChat(chat);
  }, []);

  const messages = [
    {
      id: "1",
      text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit ut aliquam",
      isMine: false,
    },
    {
      id: "2",
      text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit ut aliquam",
      isMine: true,
    },
    {
      id: "3",
      text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit ut aliquam",
      isMine: false,
    },
    {
      id: "4",
      text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit ut aliquam",
      isMine: true,
    },
  ];

  const renderMessage = ({ item }: { item: any }) => (
    <View
      style={[
        styles.messageBubble,
        item.isMine ? styles.myMessage : styles.theirMessage,
      ]}
    >
      <Text style={styles.messageText}>{item.text}</Text>
    </View>
  );

  const handleSendMessage = async () => {
    if (message.trim()) {
      try {
        await axios.post(
          `${chatServiceHost}/sendMessage`,
          { chatId, message },
          { headers: { Authorization: `Bearer ${token}` } }
        );
      } catch (error) {
        console.log(error);
      } finally {
        setMessage("");
      }
    }
  };

  useEffect(() => {
    const keyboardDidShow = () => setKeyboardVisible(true);
    const keyboardDidHide = () => setKeyboardVisible(false);

    const showSubscription = Keyboard.addListener(
      "keyboardDidShow",
      keyboardDidShow
    );
    const hideSubscription = Keyboard.addListener(
      "keyboardDidHide",
      keyboardDidHide
    );

    return () => {
      showSubscription.remove();
      hideSubscription.remove();
    };
  }, []);

  useEffect(() => {
    if (keyboardVisible || messages.length) {
      flatListRef.current?.scrollToEnd({ animated: true });
    }
  }, [keyboardVisible, messages]);

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <View style={styles.backgroundContainer}>
        <ImageBackground
          source={require("../../assets/images/chat_bg.png")}
          style={styles.background}
          resizeMode="cover"
          blurRadius={1.1}
          imageStyle={{ opacity: 0.5 }}
        />
      </View>
      <View style={styles.foreground}>
        {/* Header */}
        <View style={styles.header}>
          <TouchableOpacity onPress={router.back}>
            <Feather name="arrow-left" size={24} color="black" />
          </TouchableOpacity>
          <View style={styles.headerInfo}>
            <Text style={styles.headerTitle}>{currentChat?.chat_name}</Text>
          </View>
          <Link href={`/chats/participants/${chatId}/`}>
            <Ionicons name="people" size={24} color="black" />
          </Link>
        </View>

        {/* Chat Messages */}
        <FlatList
          ref={flatListRef}
          data={messages}
          keyExtractor={item => item.id}
          renderItem={renderMessage}
          contentContainerStyle={styles.messagesContainer}
          showsVerticalScrollIndicator={false}
        />

        {/* Message Input */}

        <View style={styles.inputContainer}>
          <TextInput
            style={styles.input}
            placeholder="Type something.."
            value={message}
            onChangeText={setMessage}
            onFocus={() => flatListRef.current?.scrollToEnd({ animated: true })}
          />
          <TouchableOpacity onPress={handleSendMessage}>
            <Ionicons name="send" size={24} color="#2b59c3" />
          </TouchableOpacity>
        </View>
      </View>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  foreground: {
    flex: 1,
    backgroundColor: "transparent", // Keeps it transparent over the background
  },
  backgroundContainer: {
    position: "absolute",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  background: {
    flex: 1,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    paddingTop: 40,
    paddingBottom: 10,
    paddingInline: 18,
    backgroundColor: "white",
    borderBottomWidth: 1,
    borderBottomColor: "#f0f0f0",
  },
  headerInfo: {
    flex: 1,
    marginLeft: 12,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: "bold",
  },
  headerStatus: {
    fontSize: 14,
    color: "green",
  },
  messagesContainer: {
    padding: 16,
    paddingBottom: 20, // Leave space for input
  },
  messageBubble: {
    padding: 12,
    borderRadius: 16,
    marginBottom: 8,
    maxWidth: "70%",
  },
  myMessage: {
    backgroundColor: "#2b59c3",
    alignSelf: "flex-end",
  },
  theirMessage: {
    backgroundColor: "white",
    alignSelf: "flex-start",
    borderWidth: 1,
    borderColor: "#ddd",
  },
  messageText: {
    color: "white",
    fontSize: 14,
  },
  inputContainer: {
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: 16,
    paddingVertical: 8,
    paddingBottom: 30,
    backgroundColor: "white",
    borderTopWidth: 1,
    borderTopColor: "#f0f0f0",
  },
  input: {
    flex: 1,
    height: 40,
    marginHorizontal: 8,
    paddingHorizontal: 12,
    backgroundColor: "#f7f7f7",
    borderRadius: 20,
  },
});

export default ChatScreen;
