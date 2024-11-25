import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
} from "react-native";
import { Feather } from "@expo/vector-icons";
import { showMessage } from "react-native-flash-message";
import { goBack } from "expo-router/build/global-state/routing";
import { router } from "expo-router";

const NewChatScreen = () => {
  const [chatName, setChatName] = useState("");

  const handleCreateChat = () => {
    if (!chatName.trim()) {
      showMessage({
        message: "Chat name cannot be empty",
        color: "red",
        backgroundColor: "#007AFF",
      });
      return;
    }

    showMessage({
      message: "Chat Created",
      description: `Chat "${chatName}" has been created.`,
      color: "#fff",
      backgroundColor: "#007AFF",
    });
    setChatName("");
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
