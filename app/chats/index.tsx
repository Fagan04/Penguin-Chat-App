import React, { useContext, useState } from "react";
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  Image,
  Alert,
} from "react-native";
import { Feather } from "@expo/vector-icons"; // Or 'react-native-vector-icons/Feather'
import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import { Link, router, useNavigation } from "expo-router";
// import Cookies from "@react-native-cookies/cookies";
import { useGlobalContext } from "@/context/GlobalProvider";

interface Chat {
  id: string;
  name: string;
  message: string;
  avatar: string;
  unread: boolean;
}

const mockChats: Chat[] = [
  {
    id: "1",
    name: "Adin Ross",
    message: "typing..",
    avatar: "https://via.placeholder.com/50",
    unread: true,
  },
  {
    id: "2",
    name: "Kamala Harris",
    message: "Ain't no party like Diddy Party..",
    avatar: "https://via.placeholder.com/50",
    unread: false,
  },
  {
    id: "3",
    name: "Donald Duck",
    message: "1 new message",
    avatar: "https://via.placeholder.com/50",
    unread: true,
  },
];

const ChatListScreen = () => {
  const [filteredChats, setFilteredChats] = useState<Chat[]>(mockChats);
  const [chats, setChats] = useState<Chat[]>(mockChats);

  const context = useGlobalContext();
  if (context == undefined) {
    throw new Error("Context not defined");
  }

  const { user, isLoggedIn, setUser, setIsLoggedIn } = context;

  function handleChangeSearch(value: string) {
    let copyChats = [...chats];
    copyChats = copyChats.filter(c =>
      c.name.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredChats(copyChats);
  }

  const handleLogout = async () => {
    try {
      //   await Cookies.clearAll();

      setUser({ id: null, name: null, email: null });
      setIsLoggedIn(false);

      Alert.alert("Logged Out", "You have been logged out successfully.");
    } catch (error) {
      console.error("Error during logout: ", error);
      Alert.alert("Error", "Failed to log out. Please try again.");
    }
  };

  const renderChatItem = ({ item }: { item: Chat }) => (
    <TouchableOpacity
      onPress={() => router.push("/chats/5")}
      style={styles.chatItem}
    >
      <Image source={{ uri: item.avatar }} style={styles.avatar} />
      <View style={styles.chatDetails}>
        <Text style={styles.chatName}>{item.name}</Text>
        <Text
          style={item.unread ? styles.chatMessageUnread : styles.chatMessage}
        >
          {item.message}
        </Text>
      </View>
      {item.unread && <View style={styles.unreadBadge} />}
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
        keyExtractor={item => item.id}
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
