export interface CookModeState {
  currentStep: number;
  totalSteps: number;
  isActive: boolean;
}

export interface TimerState {
  seconds: number;
  isRunning: boolean;
  isPaused: boolean;
  targetSeconds: number;
}

export interface ServingScale {
  original: number;
  current: number;
  multiplier: number;
}
