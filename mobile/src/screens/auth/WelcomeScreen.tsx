import React from 'react';
import {Pressable, StatusBar, StyleSheet, Text, View} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {authTheme as theme} from '../../theme/authTheme';

type WelcomeScreenProps = {onBackToLogin: () => void};

export default function WelcomeScreen({onBackToLogin}: WelcomeScreenProps) {
  return (
    <SafeAreaView style={styles.safe}>
      <StatusBar barStyle="light-content" />
      <View style={styles.content}>
        <View style={styles.mark}><Text style={styles.markText}>✦</Text></View>
        <Text style={styles.hello}>Hello World!</Text>
        <Text style={styles.title}>Welcome to Vajra</Text>
        <Text style={styles.subtitle}>Your journey to becoming better starts here.</Text>
        <Pressable accessibilityRole="button" onPress={onBackToLogin} style={({pressed}) => [styles.backButton, pressed && styles.backButtonPressed]}>
          <Text style={styles.backButtonText}>Back to Login</Text>
        </Pressable>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: {flex: 1, backgroundColor: theme.colors.background},
  content: {flex: 1, alignItems: 'center', justifyContent: 'center', paddingHorizontal: 28},
  mark: {width: 62, height: 62, borderRadius: 21, alignItems: 'center', justifyContent: 'center', backgroundColor: theme.colors.accent, marginBottom: 27},
  markText: {fontSize: 30, color: '#171321', fontWeight: '900'},
  hello: {color: theme.colors.accent, fontSize: 16, fontWeight: '700', marginBottom: 10},
  title: {color: theme.colors.text, fontSize: 30, lineHeight: 37, fontWeight: '800', textAlign: 'center'},
  subtitle: {maxWidth: 290, color: theme.colors.muted, fontSize: 15, lineHeight: 23, textAlign: 'center', marginTop: 11},
  backButton: {minHeight: 52, minWidth: 190, paddingHorizontal: 24, borderRadius: theme.radius.button, backgroundColor: theme.colors.accent, alignItems: 'center', justifyContent: 'center', marginTop: 36},
  backButtonPressed: {backgroundColor: theme.colors.accentPressed},
  backButtonText: {color: '#171321', fontSize: 14, fontWeight: '800'},
});
