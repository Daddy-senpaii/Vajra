import React, {useState} from 'react';
import {
  Alert,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StatusBar,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {AuthButton, AuthDivider, AuthField, Brand, GoogleButton} from '../../components/auth/AuthControls';
import {mockGoogleAuth, mockSignUp} from '../../services/mockAuth';
import {authTheme as theme} from '../../theme/authTheme';
import WelcomeScreen from './WelcomeScreen';

type Page = 'login' | 'signup';

export default function AuthFlow() {
  const [page, setPage] = useState<Page>('login');
  const [screen, setScreen] = useState<'auth' | 'welcome'>('auth');
  const [busy, setBusy] = useState(false);
  const [notice, setNotice] = useState('');
  const [login, setLogin] = useState({identity: '', password: ''});
  const [signup, setSignup] = useState({username: '', email: '', password: '', confirm: ''});
  const [accepted, setAccepted] = useState(false);
  const [visible, setVisible] = useState({password: false, confirm: false});
  const [errors, setErrors] = useState<Record<string, string>>({});

  const switchPage = (next: Page) => {
    setPage(next);
    setErrors({});
    setNotice('');
  };

  const submit = async () => {
    const nextErrors: Record<string, string> = {};
    if (page === 'login') {
      if (!login.identity.trim()) nextErrors.identity = 'Enter your email or username.';
      if (!login.password) nextErrors.password = 'Enter your password.';
    } else {
      if (!/^[a-zA-Z0-9_.]{3,20}$/.test(signup.username)) nextErrors.username = 'Use 3–20 letters, numbers, dots or underscores.';
      if (!/^\S+@\S+\.\S+$/.test(signup.email.trim())) nextErrors.email = 'Enter a valid email address.';
      if (signup.password.length < 8) nextErrors.signupPassword = 'Use at least 8 characters.';
      if (signup.confirm !== signup.password) nextErrors.confirm = 'Passwords do not match.';
      if (!accepted) nextErrors.terms = 'Please accept the terms to continue.';
    }
    setErrors(nextErrors);
    setNotice('');
    if (Object.keys(nextErrors).length) return;
    if (page === 'login') {
      setScreen('welcome');
      return;
    }
    setBusy(true);
    try {
      const result = await mockSignUp();
      setNotice(result.message);
    } finally {
      setBusy(false);
    }
  };

  const google = async () => {
    setBusy(true);
    setNotice('');
    try {
      const result = await mockGoogleAuth();
      setNotice(result.message);
    } finally {
      setBusy(false);
    }
  };

  if (screen === 'welcome') {
    return <WelcomeScreen onBackToLogin={() => {setScreen('auth'); setPage('login'); setNotice(''); setErrors({});}} />;
  }

  return (
    <SafeAreaView style={styles.safe}>
      <StatusBar barStyle="light-content" />
      <KeyboardAvoidingView style={styles.flex} behavior={Platform.OS === 'ios' ? 'padding' : undefined}>
        <ScrollView contentContainerStyle={styles.scroll} keyboardShouldPersistTaps="handled" showsVerticalScrollIndicator={false}>
          <View style={styles.content}>
            <Brand />
            <Text style={styles.heading}>{page === 'login' ? 'Welcome back' : 'Create account'}</Text>
            <Text style={styles.subtitle}>{page === 'login' ? 'Sign in and keep building the habits that matter.' : 'Start making steady progress, one day at a time.'}</Text>

            {page === 'login' ? (
              <>
                <AuthField label="Email or username" placeholder="you@example.com" value={login.identity} onChangeText={identity => setLogin({...login, identity})} keyboardType="email-address" returnKeyType="next" error={errors.identity} />
                <AuthField label="Password" placeholder="Enter your password" value={login.password} onChangeText={password => setLogin({...login, password})} secureTextEntry={!visible.password} returnKeyType="done" error={errors.password} rightAdornment={<VisibilityToggle visible={visible.password} onPress={() => setVisible({...visible, password: !visible.password})} />} />
                <Pressable accessibilityRole="button" onPress={() => Alert.alert('Forgot password', 'Password recovery will be available when authentication is connected.')} style={styles.forgot}><Text style={styles.link}>Forgot password?</Text></Pressable>
              </>
            ) : (
              <>
                <AuthField label="Username" placeholder="Choose a username" value={signup.username} onChangeText={username => setSignup({...signup, username})} error={errors.username} />
                <AuthField label="Email" placeholder="you@example.com" value={signup.email} onChangeText={email => setSignup({...signup, email})} keyboardType="email-address" error={errors.email} />
                <AuthField label="Password" placeholder="At least 8 characters" value={signup.password} onChangeText={password => setSignup({...signup, password})} secureTextEntry={!visible.password} error={errors.signupPassword} rightAdornment={<VisibilityToggle visible={visible.password} onPress={() => setVisible({...visible, password: !visible.password})} />} />
                <AuthField label="Confirm password" placeholder="Re-enter your password" value={signup.confirm} onChangeText={confirm => setSignup({...signup, confirm})} secureTextEntry={!visible.confirm} error={errors.confirm} rightAdornment={<VisibilityToggle visible={visible.confirm} onPress={() => setVisible({...visible, confirm: !visible.confirm})} />} />
                <Pressable accessibilityRole="checkbox" accessibilityState={{checked: accepted}} onPress={() => setAccepted(!accepted)} style={styles.termsRow}>
                  <View style={[styles.checkbox, accepted && styles.checkboxChecked]}>{accepted ? <Text style={styles.check}>✓</Text> : null}</View>
                  <Text style={styles.termsText}>I agree to the <Text style={styles.link}>Terms</Text> and <Text style={styles.link}>Privacy Policy</Text>.</Text>
                </Pressable>
                {errors.terms ? <Text style={styles.formError}>{errors.terms}</Text> : null}
              </>
            )}

            {notice ? <Text accessibilityRole="alert" style={styles.notice}>{notice}</Text> : null}
            <AuthButton title={page === 'login' ? 'Sign In' : 'Create Account'} onPress={submit} disabled={busy} loading={busy} />
            <AuthDivider />
            <GoogleButton title={page === 'login' ? 'Sign in with Google' : 'Sign up with Google'} onPress={google} disabled={busy} />
            <View style={styles.footer}>
              <Text style={styles.footerText}>{page === 'login' ? "Don't have an account?" : 'Already have an account?'}</Text>
              <Pressable accessibilityRole="button" onPress={() => switchPage(page === 'login' ? 'signup' : 'login')}><Text style={styles.link}>{page === 'login' ? ' Sign Up' : ' Sign In'}</Text></Pressable>
            </View>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

function VisibilityToggle({visible, onPress}: {visible: boolean; onPress: () => void}) {
  return <Pressable accessibilityRole="button" accessibilityLabel={visible ? 'Hide password' : 'Show password'} onPress={onPress} hitSlop={10}><Text style={styles.showHide}>{visible ? 'Hide' : 'Show'}</Text></Pressable>;
}

const styles = StyleSheet.create({
  flex: {flex: 1},
  safe: {flex: 1, backgroundColor: theme.colors.background},
  scroll: {flexGrow: 1, justifyContent: 'center', paddingHorizontal: 24, paddingVertical: 30},
  content: {width: '100%', maxWidth: 440, alignSelf: 'center'},
  heading: {color: theme.colors.text, fontSize: 30, lineHeight: 36, fontWeight: '800', letterSpacing: -0.7},
  subtitle: {color: theme.colors.muted, fontSize: 14, lineHeight: 21, marginTop: 7, marginBottom: 27},
  forgot: {alignSelf: 'flex-end', marginTop: -3, marginBottom: 17, paddingVertical: 4},
  link: {color: theme.colors.accent, fontSize: 13, fontWeight: '700'},
  showHide: {color: theme.colors.accent, fontSize: 12, fontWeight: '700'},
  termsRow: {flexDirection: 'row', alignItems: 'center', marginTop: 0, marginBottom: 16, minHeight: 32},
  checkbox: {width: 19, height: 19, borderRadius: 5, borderWidth: 1, borderColor: theme.colors.subtle, alignItems: 'center', justifyContent: 'center', marginRight: 10},
  checkboxChecked: {backgroundColor: theme.colors.accent, borderColor: theme.colors.accent},
  check: {color: '#171321', fontWeight: '900', fontSize: 13, marginTop: -1},
  termsText: {color: theme.colors.muted, fontSize: 12, flexShrink: 1},
  formError: {color: theme.colors.error, fontSize: 11, marginTop: -10, marginBottom: 10},
  notice: {color: theme.colors.positive, fontSize: 12, lineHeight: 18, marginTop: 1, marginBottom: 7},
  footer: {flexDirection: 'row', justifyContent: 'center', alignItems: 'center', marginTop: 25, minHeight: 28},
  footerText: {color: theme.colors.muted, fontSize: 13},
});
