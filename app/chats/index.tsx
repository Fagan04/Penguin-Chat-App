import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  Image,
} from "react-native";
import { Feather } from "@expo/vector-icons"; // Or 'react-native-vector-icons/Feather'
import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import { Link, router } from "expo-router";
import { useGlobalContext } from "@/context/GlobalProvider";
import showSuccessMessage from "@/utils/showSuccessMessage";
import axios from "axios";
import { chatServiceHost } from "@/constants/backendUrl";
import { Chat } from "@/types/Chat";
import { io } from "socket.io-client";

const ChatListScreen = () => {
  const [filteredChats, setFilteredChats] = useState<Chat[]>([]);

  const context = useGlobalContext();
  if (context == undefined) throw new Error("Context not defined");

  const { token, setToken, chats, setChats, setConnection, connection } =
    context;

  useEffect(() => {
    if (!token) router.replace("/");
  }, [token]);

  useEffect(() => {
    try {
      const getData = async () => {
        const { data }: { data: Chat[] | { message: string } } =
          await axios.get(`${chatServiceHost}/accessChat`, {
            headers: { Authorization: `Bearer ${token}` },
          });
        if (Array.isArray(data)) {
          setChats(data);
          setFilteredChats(data);
        }
      };
      getData();
    } catch (error) {
      console.log(error);
    }
  }, []);

  // todo: WEBSOCKET

  useEffect(() => {
    try {
      const ws = new WebSocket("ws://192.168.31.208:8081/ws?chatID=1", [
        "Authorization",
        `Bearer ${token}`,
      ]);

      ws.addEventListener("error", ev => {
        console.error(ev);
      });

      ws.addEventListener("text", () => {
        console.log("NEW TEXT WEBSOCKET");
      });

      ws.addEventListener("open", () => {
        setConnection(ws);
        console.log("OPEN");
      });

      showSuccessMessage("Successfull connection");
      return ws.close;
    } catch (error) {
      console.log("Websoket error", error);
      showSuccessMessage("Failed to connect", false);
    }
  }, []);

  function handleChangeSearch(value: string) {
    let copyChats = [...chats];
    copyChats = copyChats.filter(c =>
      c.chat_name.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredChats(copyChats);
  }

  const handleLogout = async () => {
    try {
      setToken("");
      showSuccessMessage(
        "You have been logged out successfully.",
        true,
        "Logged out"
      );
      router.replace("/");
    } catch (error) {
      console.error("Error during logout: ", error);
      showSuccessMessage("Oops. An error occured");
    }
  };

  const renderChatItem = ({ item }: { item: Chat }) => (
    <TouchableOpacity
      onPress={() => router.push(`/chats/${item.id}`)}
      style={styles.chatItem}
    >
      <Image
        source={{ uri: "https://via.placeholder.com/50" }}
        style={styles.avatar}
      />
      <View style={styles.chatDetails}>
        <Text style={styles.chatName}>{item.chat_name}</Text>
        <Text style={true ? styles.chatMessageUnread : styles.chatMessage}>
          {"salam"}
        </Text>
      </View>
      {true && <View style={styles.unreadBadge} />}
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={handleLogout}>
          <MaterialIcons name="logout" size={24} color="black" />
        </TouchableOpacity>
        <Link href="/login">
          <Image
            source={{ uri: "https://via.placeholder.com/50" }} // Replace with a user avatar image
            style={styles.profilePicture}
          />
        </Link>
      </View>

      {/* Chats Section */}
      <Text style={styles.sectionTitle}>Chats</Text>
      <View style={styles.searchContainer}>
        <TextInput
          onChangeText={value => handleChangeSearch(value)}
          placeholder="search conversation..."
          style={styles.searchInput}
        />
        <Feather name="message-square" size={20} color="gray" />
      </View>

      {/* Chat List */}
      <FlatList
        data={filteredChats}
        keyExtractor={item => item.id.toString()}
        renderItem={renderChatItem}
        contentContainerStyle={styles.chatList}
      />

      {/* Floating Action Button */}

      <Link href="/chats/new" style={styles.fabLink}>
        <View style={styles.fab}>
          <Feather name="message-square" size={24} color="white" />
          <Text style={styles.fabText}>New Chat</Text>
        </View>
      </Link>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingTop: 30,
    backgroundColor: "#dce6ff",
    paddingHorizontal: 16,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    marginVertical: 16,
  },
  profilePicture: {
    width: 40,
    height: 40,
    borderRadius: 20,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: "bold",
    color: "#2b59c3",
    marginVertical: 16,
  },
  searchContainer: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "white",
    borderRadius: 10,
    paddingHorizontal: 12,
    paddingVertical: 8,
    marginBottom: 16,
  },
  searchInput: {
    flex: 1,
    fontSize: 16,
    marginRight: 8,
  },
  chatList: {
    paddingBottom: 80,
  },
  chatItem: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "white",
    borderRadius: 10,
    padding: 12,
    marginBottom: 8,
  },
  avatar: {
    width: 50,
    height: 50,
    borderRadius: 25,
    marginRight: 12,
  },
  chatDetails: {
    flex: 1,
  },
  chatName: {
    fontSize: 16,
    fontWeight: "bold",
    marginBottom: 4,
  },
  chatMessage: {
    fontSize: 14,
    color: "gray",
  },
  chatMessageUnread: {
    fontSize: 14,
    color: "#2b59c3",
    fontWeight: "bold",
  },
  unreadBadge: {
    width: 10,
    height: 10,
    backgroundColor: "#2b59c3",
    borderRadius: 5,
  },
  fabLink: {
    position: "absolute",
    right: 15,
    bottom: 30,
  },
  fab: {
    display: "flex",
    flexDirection: "row",
    backgroundColor: "#2b59c3",
    alignItems: "center",
    paddingHorizontal: 16,
    paddingVertical: 10,
    borderRadius: 25,
    elevation: 5,
  },
  fabText: {
    color: "white",
    fontSize: 16,
    marginLeft: 8,
  },
});

export default ChatListScreen;
