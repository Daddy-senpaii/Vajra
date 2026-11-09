import React from 'react';
import {
  ActivityIndicator,
  Pressable,
  StyleSheet,
  Text,
  TextInput,
  TextInputProps,
  View,
} from 'react-native';
import {authTheme as theme} from '../../theme/authTheme';

export function AuthField({label, error, rightAdornment, ...props}: TextInputProps & {label: string; error?: string; rightAdornment?: React.ReactNode}) {
  return (
    <View style={styles.fieldWrap}>
      <Text style={styles.label}>{label}</Text>
      <View style={[styles.inputRow, error ? styles.inputError : null]}>
        <TextInput
          placeholderTextColor={theme.colors.subtle}
          selectionColor={theme.colors.accent}
          autoCapitalize="none"
          style={styles.input}
          accessibilityLabel={label}
          {...props}
        />
        {rightAdornment}
      </View>
      {error ? <Text style={styles.errorText}>{error}</Text> : null}
    </View>
  );
}

export function AuthButton({title, onPress, disabled, loading}: {
  title: string; onPress: () => void; disabled?: boolean; loading?: boolean;
}) {
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityState={{disabled: Boolean(disabled || loading), busy: loading}}
      onPress={onPress}
      disabled={disabled || loading}
      style={({pressed}) => [styles.primary, (disabled || loading) && styles.disabled, pressed && !disabled && styles.pressed]}>
      {loading ? <ActivityIndicator color="#171321" /> : <Text style={styles.primaryText}>{title}</Text>}
    </Pressable>
  );
}

export function GoogleButton({title, onPress, disabled}: {title: string; onPress: () => void; disabled?: boolean}) {
  return (
    <Pressable accessibilityRole="button" accessibilityState={{disabled: Boolean(disabled)}} onPress={onPress} disabled={disabled} style={styles.googleButton}>
      <Text accessibilityLabel="Google" style={styles.googleG}>G</Text>
      <Text style={styles.googleText}>{title}</Text>
    </Pressable>
  );
}

export function AuthDivider() {
  return <View style={styles.dividerRow}><View style={styles.divider} /><Text style={styles.or}>OR</Text><View style={styles.divider} /></View>;
}

export function Brand() {
  return <View style={styles.brandRow}><View style={styles.brandMark}><Text style={styles.brandGlyph}>✦</Text></View><Text style={styles.brandName}>vajra</Text></View>;
}

const styles = StyleSheet.create({
  fieldWrap: {marginBottom: 15},
  label: {color: theme.colors.muted, fontSize: 12, fontWeight: '600', marginBottom: 8},
  inputRow: {minHeight: 53, borderRadius: theme.radius.input, borderWidth: 1, borderColor: theme.colors.border, backgroundColor: theme.colors.surface, flexDirection: 'row', alignItems: 'center', paddingRight: 14},
  input: {flex: 1, minHeight: 51, color: theme.colors.text, paddingHorizontal: 15, fontSize: 15},
  inputError: {borderColor: theme.colors.error},
  errorText: {color: theme.colors.error, fontSize: 11, marginTop: 5},
  primary: {minHeight: 54, borderRadius: theme.radius.button, backgroundColor: theme.colors.accent, alignItems: 'center', justifyContent: 'center', marginTop: 7},
  primaryText: {color: '#171321', fontSize: 15, fontWeight: '800'},
  pressed: {backgroundColor: theme.colors.accentPressed},
  disabled: {opacity: 0.55},
  googleButton: {height: 53, borderRadius: theme.radius.button, borderWidth: 1, borderColor: theme.colors.border, backgroundColor: theme.colors.surface, alignItems: 'center', justifyContent: 'center', flexDirection: 'row'},
  googleG: {fontSize: 20, fontWeight: '800', color: '#4285F4', marginRight: 10},
  googleText: {color: theme.colors.text, fontSize: 14, fontWeight: '600'},
  dividerRow: {flexDirection: 'row', alignItems: 'center', gap: 13, marginVertical: 20},
  divider: {height: 1, backgroundColor: theme.colors.border, flex: 1},
  or: {color: theme.colors.subtle, fontWeight: '700', fontSize: 10, letterSpacing: 1.5},
  brandRow: {flexDirection: 'row', alignItems: 'center', gap: 9, marginBottom: 31},
  brandMark: {width: 31, height: 31, borderRadius: 11, backgroundColor: theme.colors.accent, alignItems: 'center', justifyContent: 'center'},
  brandGlyph: {color: '#171321', fontSize: 17, fontWeight: '900'},
  brandName: {color: theme.colors.text, fontSize: 21, letterSpacing: 0.4, fontWeight: '800'},
});
