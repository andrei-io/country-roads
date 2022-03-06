import { Inter_500Medium, Inter_600SemiBold, useFonts } from '@expo-google-fonts/inter';
import { Router } from './router/router';

export default function App() {
  let [fontsLoaded] = useFonts({
    Inter_600SemiBold,
    Inter_500Medium,
  });

  if (!fontsLoaded) {
    return null;
  }
  return <Router />;
}
