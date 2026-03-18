import { CookModeRoot } from './CookModeRoot';
import { CookModeStep } from './CookModeStep';
import { CookModeTimer } from './CookModeTimer';
import { CookModeControls } from './CookModeControls';

export const CookMode = {
  Root: CookModeRoot,
  Step: CookModeStep,
  Timer: CookModeTimer,
  Controls: CookModeControls,
};

export { useCookMode } from './CookModeRoot';
