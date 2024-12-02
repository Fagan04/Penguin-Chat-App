import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
} from "react-native";
import { Feather } from "@expo/vector-icons";
import { router } from "expo-router";
import showSuccessMessage from "@/utils/showSuccessMessage";
import axios from "axios";
import { useGlobalContext } from "@/context/GlobalProvider";
import { chatServiceHost } from "@/constants/backendUrl";

const NewChatScreen = () => {
  const [chatName, setChatName] = useState("");

  const context = useGlobalContext();

  if (!context) throw new Error("Context not defined");

  const { token } = context;

  const handleCreateChat = async () => {
    if (!chatName.trim()) {
      showSuccessMessage("Chat name cannot be empty", false);
      return;
    }

    try {
      const { data } = await axios.post(
        `${chatServiceHost}/createChat`,
        {
          chat_name: chatName,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      console.log(data);
      showSuccessMessage(`New chat named: ${chatName} created.`);
      setChatName("");
      router.replace("/chats");
    } catch (error) {
      console.log(error);
      if (axios.isAxiosError(error)) {
        showSuccessMessage(error.response?.data.trim(), false);
      }
    }
  };
  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={router.back}>
          <Feather name="arrow-left" size={24} color="black" />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>New Chat</Text>
      </View>

      {/* Input Field */}
      <View style={styles.inputContainer}>
        <Text style={styles.label}>Enter Chat Name</Text>
        <TextInput
          style={styles.input}
          placeholder="Chat Name"
          value={chatName}
          onChangeText={setChatName}
        />
      </View>

      {/* Create Button */}
      <TouchableOpacity style={styles.createButton} onPress={handleCreateChat}>
        <Text style={styles.createButtonText}>Create Chat</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#dce6ff",
    paddingHorizontal: 16,
    paddingTop: 46,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 24,
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: "bold",
    marginLeft: 12,
  },
  inputContainer: {
    marginBottom: 24,
  },
  label: {
    fontSize: 16,
    fontWeight: "bold",
    color: "#2b59c3",
    marginBottom: 8,
  },
  input: {
    backgroundColor: "white",
    borderRadius: 10,
    padding: 12,
    fontSize: 16,
    borderWidth: 1,
    borderColor: "#d0d4e2",
  },
  createButton: {
    backgroundColor: "#2b59c3",
    paddingVertical: 12,
    borderRadius: 10,
    alignItems: "center",
    elevation: 3,
  },
  createButtonText: {
    color: "white",
    fontSize: 16,
    fontWeight: "bold",
  },
});

export default NewChatScreen;
