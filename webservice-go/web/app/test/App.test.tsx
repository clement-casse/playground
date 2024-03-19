import renderer from 'react-test-renderer';
import { test } from 'vitest'
import { NavItems } from '../src/App';

test('nav items contains some links', () => {
  const component = renderer.create(<NavItems />)

  component.getInstance
})