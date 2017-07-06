import { animate, AnimationEntryMetadata, state, style, transition, trigger } from '@angular/core';

export const fadeInAnimation: AnimationEntryMetadata =
  trigger('fadeInAnimation', [
    state('inactive', style({
      opacity: 0,
    })),
    state('active', style({
      opacity: 1,
    })),
    transition('inactive => active', animate('100ms ease-in')),
    transition('active => inactive', animate('100ms ease-out')),
  ]);
