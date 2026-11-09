import {useState} from 'react';
import {
  Pressable,
  ScrollView,
  StatusBar,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import {SafeAreaProvider, SafeAreaView} from 'react-native-safe-area-context';
import AuthFlow from './src/screens/auth/AuthFlow';

const initialTasks = [
  {id: 'move', icon: '↗', title: 'Move your body', detail: '30 min movement', color: '#DDF06B'},
  {id: 'read', icon: '▤', title: 'Read a little', detail: '10 pages', color: '#D7D2FF'},
  {id: 'water', icon: '◒', title: 'Drink more water', detail: '6 of 8 glasses', color: '#BCE9E6'},
];

export function HomePreview() {
  const [completed, setCompleted] = useState<string[]>(['water']);
  const doneCount = completed.length;

  const toggleTask = (id: string) => {
    setCompleted(current =>
      current.includes(id) ? current.filter(task => task !== id) : [...current, id],
    );
  };

  return (
    <SafeAreaProvider>
      <StatusBar barStyle="dark-content" />
      <SafeAreaView style={styles.safeArea}>
        <View style={styles.screen}>
          <ScrollView
            contentContainerStyle={styles.scrollContent}
            showsVerticalScrollIndicator={false}>
            <View style={styles.header}>
              <View>
                <Text style={styles.eyebrow}>MONDAY, JUNE 16</Text>
                <Text style={styles.greeting}>Good morning, Alex <Text style={styles.wave}>✳</Text></Text>
              </View>
              <Pressable accessibilityRole="button" accessibilityLabel="Notifications" style={styles.bellButton}>
                <Text style={styles.bell}>♧</Text>
                <View style={styles.notificationDot} />
              </Pressable>
            </View>

            <View style={styles.streakRow}>
              <View style={styles.streakIcon}><Text style={styles.flame}>✦</Text></View>
              <Text style={styles.streakText}><Text style={styles.streakNumber}>7 day streak</Text>  ·  Keep showing up</Text>
            </View>

            <View style={styles.heroCard}>
              <View style={styles.heroTopline}>
                <Text style={styles.heroLabel}>YOUR ACTIVE CHALLENGE</Text>
                <View style={styles.activeBadge}><View style={styles.activeDot} /><Text style={styles.activeText}>IN PROGRESS</Text></View>
              </View>
              <Text style={styles.challengeTitle}>A little better,{"\n"}every day.</Text>
              <Text style={styles.challengeSubtitle}>Build a routine that feels like you.</Text>
              <View style={styles.progressMeta}>
                <Text style={styles.progressLabel}>DAY 12 <Text style={styles.progressMuted}>OF 21</Text></Text>
                <Text style={styles.progressPercent}>57%</Text>
              </View>
              <View style={styles.progressTrack}><View style={styles.progressFill} /></View>
              <View style={styles.heroFooter}>
                <View style={styles.avatarStack}>
                  <View style={[styles.avatar, styles.avatarOne]}><Text style={styles.avatarText}>J</Text></View>
                  <View style={[styles.avatar, styles.avatarTwo]}><Text style={styles.avatarText}>M</Text></View>
                  <View style={[styles.avatar, styles.avatarThree]}><Text style={styles.avatarText}>+</Text></View>
                </View>
                <Text style={styles.withFriends}>You and 2 friends</Text>
                <Text style={styles.arrow}>↗</Text>
              </View>
              <View style={styles.heroOrb} />
              <View style={styles.heroOrbSmall} />
            </View>

            <View style={styles.sectionHeader}>
              <View>
                <Text style={styles.sectionTitle}>Today, your way</Text>
                <Text style={styles.sectionSubtitle}>{doneCount} of {initialTasks.length} habits complete</Text>
              </View>
              <Pressable accessibilityRole="button"><Text style={styles.seeAll}>SEE ALL</Text></Pressable>
            </View>

            <View style={styles.taskList}>
              {initialTasks.map(task => {
                const isDone = completed.includes(task.id);
                return (
                  <Pressable
                    key={task.id}
                    accessibilityRole="checkbox"
                    accessibilityState={{checked: isDone}}
                    onPress={() => toggleTask(task.id)}
                    style={[styles.taskCard, isDone && styles.taskCardDone]}>
                    <View style={[styles.taskIcon, {backgroundColor: task.color}]}><Text style={styles.taskGlyph}>{task.icon}</Text></View>
                    <View style={styles.taskCopy}>
                      <Text style={[styles.taskTitle, isDone && styles.taskTitleDone]}>{task.title}</Text>
                      <Text style={styles.taskDetail}>{task.detail}</Text>
                    </View>
                    <View style={[styles.checkCircle, isDone && styles.checkCircleDone]}>
                      {isDone ? <Text style={styles.checkMark}>✓</Text> : null}
                    </View>
                  </Pressable>
                );
              })}
            </View>

            <View style={styles.quoteCard}>
              <View style={styles.quoteMark}><Text style={styles.quoteMarkText}>“</Text></View>
              <View style={styles.quoteCopy}>
                <Text style={styles.quoteText}>Small steps still move you forward.</Text>
                <Text style={styles.quoteByline}>A NOTE FOR TODAY</Text>
              </View>
              <Text style={styles.sparkle}>✳</Text>
            </View>
          </ScrollView>

          <View style={styles.tabBar}>
            <Pressable style={styles.tabItem}><Text style={styles.tabIconActive}>⌂</Text><Text style={styles.tabLabelActive}>Home</Text></Pressable>
            <Pressable style={styles.tabItem}><Text style={styles.tabIcon}>◉</Text><Text style={styles.tabLabel}>Friends</Text></Pressable>
            <Pressable style={styles.tabItem}><Text style={styles.tabIcon}>▥</Text><Text style={styles.tabLabel}>Progress</Text></Pressable>
            <Pressable style={styles.tabItem}><Text style={styles.tabIcon}>○</Text><Text style={styles.tabLabel}>You</Text></Pressable>
          </View>
        </View>
      </SafeAreaView>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  safeArea: {flex: 1, backgroundColor: '#F6F6F1'},
  screen: {flex: 1, backgroundColor: '#F6F6F1'},
  scrollContent: {paddingHorizontal: 22, paddingTop: 16, paddingBottom: 24},
  header: {flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', marginBottom: 19},
  eyebrow: {color: '#93968E', fontSize: 10, fontWeight: '800', letterSpacing: 1.5, marginBottom: 7},
  greeting: {color: '#20251F', fontSize: 23, fontWeight: '700', letterSpacing: -0.8},
  wave: {color: '#F2A35B', fontSize: 19},
  bellButton: {width: 44, height: 44, borderRadius: 16, backgroundColor: '#FFFFFF', alignItems: 'center', justifyContent: 'center'},
  bell: {fontSize: 23, color: '#273027', transform: [{rotate: '180deg'}]},
  notificationDot: {position: 'absolute', top: 9, right: 10, width: 7, height: 7, borderRadius: 4, backgroundColor: '#F07C62', borderWidth: 1, borderColor: '#FFFFFF'},
  streakRow: {flexDirection: 'row', alignItems: 'center', marginBottom: 19},
  streakIcon: {width: 25, height: 25, borderRadius: 9, backgroundColor: '#FFE5C8', alignItems: 'center', justifyContent: 'center', marginRight: 9},
  flame: {color: '#E88943', fontSize: 15, fontWeight: '800'},
  streakText: {color: '#858A81', fontSize: 12},
  streakNumber: {color: '#343B32', fontWeight: '700'},
  heroCard: {backgroundColor: '#222C27', borderRadius: 27, padding: 21, overflow: 'hidden', marginBottom: 28},
  heroTopline: {flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', marginBottom: 18},
  heroLabel: {color: '#A9B2A8', fontSize: 9, fontWeight: '800', letterSpacing: 1.3},
  activeBadge: {flexDirection: 'row', alignItems: 'center', borderRadius: 20, paddingVertical: 6, paddingHorizontal: 8, backgroundColor: 'rgba(255,255,255,0.09)'},
  activeDot: {width: 5, height: 5, borderRadius: 3, backgroundColor: '#DDF06B', marginRight: 5},
  activeText: {color: '#E6EBDD', fontSize: 8, fontWeight: '800', letterSpacing: 0.6},
  challengeTitle: {color: '#FFFFFF', fontSize: 29, lineHeight: 33, fontWeight: '700', letterSpacing: -0.7},
  challengeSubtitle: {color: '#A9B2A8', fontSize: 12, marginTop: 7, marginBottom: 21},
  progressMeta: {flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8},
  progressLabel: {color: '#EAF0E7', fontSize: 10, fontWeight: '800', letterSpacing: 0.8},
  progressMuted: {color: '#89938A'},
  progressPercent: {color: '#DDF06B', fontSize: 12, fontWeight: '800'},
  progressTrack: {height: 6, borderRadius: 4, backgroundColor: '#465148', overflow: 'hidden'},
  progressFill: {width: '57%', height: 6, borderRadius: 4, backgroundColor: '#DDF06B'},
  heroFooter: {flexDirection: 'row', alignItems: 'center', marginTop: 19},
  avatarStack: {flexDirection: 'row', marginRight: 9},
  avatar: {width: 23, height: 23, borderRadius: 12, alignItems: 'center', justifyContent: 'center', borderWidth: 1.5, borderColor: '#222C27', marginRight: -5},
  avatarOne: {backgroundColor: '#E7A886'},
  avatarTwo: {backgroundColor: '#ADB9E8'},
  avatarThree: {backgroundColor: '#59655A'},
  avatarText: {color: '#243027', fontSize: 9, fontWeight: '800'},
  withFriends: {color: '#B6BFB5', fontSize: 10, flex: 1},
  arrow: {color: '#DDF06B', fontSize: 18, fontWeight: '700'},
  heroOrb: {position: 'absolute', width: 110, height: 110, borderRadius: 55, borderWidth: 1, borderColor: 'rgba(221,240,107,0.12)', top: 75, right: -57},
  heroOrbSmall: {position: 'absolute', width: 64, height: 64, borderRadius: 32, borderWidth: 1, borderColor: 'rgba(221,240,107,0.12)', top: 98, right: -34},
  sectionHeader: {flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', marginBottom: 13},
  sectionTitle: {color: '#20251F', fontSize: 20, fontWeight: '700', letterSpacing: -0.5},
  sectionSubtitle: {color: '#8D9289', fontSize: 11, marginTop: 4},
  seeAll: {color: '#687A58', fontSize: 9, fontWeight: '800', letterSpacing: 1},
  taskList: {gap: 9},
  taskCard: {minHeight: 70, flexDirection: 'row', alignItems: 'center', backgroundColor: '#FFFFFF', borderRadius: 19, paddingHorizontal: 12, paddingVertical: 10},
  taskCardDone: {backgroundColor: '#FBFCF8'},
  taskIcon: {width: 45, height: 45, borderRadius: 15, alignItems: 'center', justifyContent: 'center', marginRight: 12},
  taskGlyph: {color: '#263127', fontSize: 21, fontWeight: '700'},
  taskCopy: {flex: 1},
  taskTitle: {color: '#293029', fontSize: 13, fontWeight: '700', marginBottom: 4},
  taskTitleDone: {color: '#7D8977'},
  taskDetail: {color: '#92978E', fontSize: 10},
  checkCircle: {width: 23, height: 23, borderRadius: 12, borderWidth: 1.5, borderColor: '#D8DDD4', alignItems: 'center', justifyContent: 'center'},
  checkCircleDone: {backgroundColor: '#657B57', borderColor: '#657B57'},
  checkMark: {color: '#FFFFFF', fontSize: 13, fontWeight: '800', marginTop: -1},
  quoteCard: {flexDirection: 'row', alignItems: 'center', backgroundColor: '#E9E9DF', borderRadius: 20, padding: 15, marginTop: 17},
  quoteMark: {width: 35, height: 35, borderRadius: 12, backgroundColor: '#D8DCCB', alignItems: 'center', justifyContent: 'center', marginRight: 11},
  quoteMarkText: {color: '#687A58', fontSize: 29, fontWeight: '700', marginTop: 4},
  quoteCopy: {flex: 1},
  quoteText: {color: '#485143', fontSize: 12, fontWeight: '600'},
  quoteByline: {color: '#8B9183', fontSize: 8, fontWeight: '800', letterSpacing: 1, marginTop: 5},
  sparkle: {color: '#A6B477', fontSize: 18},
  tabBar: {height: 67, flexDirection: 'row', alignItems: 'center', justifyContent: 'space-around', backgroundColor: '#FFFFFF', borderTopWidth: 1, borderTopColor: '#ECEDE7', paddingBottom: 4},
  tabItem: {minWidth: 55, alignItems: 'center', justifyContent: 'center'},
  tabIconActive: {color: '#536B47', fontSize: 22, fontWeight: '700', lineHeight: 25},
  tabIcon: {color: '#A1A69D', fontSize: 19, lineHeight: 25},
  tabLabelActive: {color: '#536B47', fontSize: 9, fontWeight: '800', marginTop: 2},
  tabLabel: {color: '#989E95', fontSize: 9, fontWeight: '600', marginTop: 2},
});

export default function App() {
  return (
    <SafeAreaProvider>
      <AuthFlow />
    </SafeAreaProvider>
  );
}
