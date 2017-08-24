import { animate, AnimationEntryMetadata, state, style, transition, trigger, keyframes } from '@angular/core';

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

export const flyInOutAnimation: AnimationEntryMetadata =
  trigger('flyInOut', [
    state('in', style({ transform: 'translateX(0)' })),
    transition('void => *', [
      style({
        opacity: 0,
        transform: 'translateX(-100%)'
      }),
      animate(250)
    ]),
    transition('* => void', [
      animate(250, style({ transform: 'translateX(100%)' }))
    ])
  ])

export const bounceInAnimation: AnimationEntryMetadata =
  trigger('bounceIn', [
    state('in', style({ opacity: 1, transform: 'scale3d(1, 1, 1)' })),
    transition('void => *', [
      animate(300, keyframes([
        style({ opacity: 0, transform: 'scale3d(.3, .3, .3)', offset: 0 }),
        style({ transform: 'scale3d(1.1, 1.1, 1.1)', offset: 0.2 }),
        style({ transform: 'scale3d(.9, .9, .9)', offset: 0.4 }),
        style({ opacity: 1, transform: 'scale3d(1.03, 1.03, 1.03)', offset: 0.6 }),
        style({ transform: 'scale3d(.97, .97, .97)', offset: 0.8 }),
        style({ opacity: 1, transform: 'scale3d(1, 1, 1)', offset: 1 })
      ]))
    ]),
    transition('* => void', [
      animate(250, style({ transform: 'translateX(100%)' }))
    ])
  ])
