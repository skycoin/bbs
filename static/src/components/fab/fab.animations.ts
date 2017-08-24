import { animate, AnimationEntryMetadata, state, style, transition, trigger } from '@angular/core';

export const flyInOutAnimation: AnimationEntryMetadata =
  trigger('flyInOut', [
    state('in', style({ transform: 'translateY(0)' })),
    transition('void => *', [
      style({
        opacity: 0,
        transform: 'translateY(15%)'
      }),
      animate(150)
    ]),
    transition('* => void', [
      animate(150, style({
        opacity: 0,
        transform: 'translateY(15%)'
      }))
    ])
  ])

export const rotate45Animation: AnimationEntryMetadata =
  trigger('rotate45', [
    state('inactive', style({
      transform: 'rotate(0)'
    })),
    state('active', style({
      transform: 'rotate(45deg)'
    })),
    transition('inactive => active', animate('100ms ease-in')),
    transition('active => inactive', animate('100ms ease-out'))
  ])
