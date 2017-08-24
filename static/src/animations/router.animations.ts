import { animate, AnimationEntryMetadata, state, style, transition, trigger } from '@angular/core';

// Component transition animations
export const slideInLeftAnimation: AnimationEntryMetadata =
  trigger('routeAnimation', [
    state('*',
      style({
        opacity: 1
      })
    ),
    transition(':enter', [
      style({
        opacity: 0
      }),
      animate('0.3s ease-in')
    ]),
    transition(':leave', [
      animate('0.3s ease-out', style({
        opacity: 0
      }))
    ])
  ]);

