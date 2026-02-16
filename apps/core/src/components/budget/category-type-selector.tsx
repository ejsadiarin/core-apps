'use client';

import { useState } from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select';
import { useUpdateCategoryType, useCategories } from '@/hooks/use-budget';
import type { Category } from '@/types/api';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface CategoryTypeSelectorProps {
  category: Category;
  onUpdate?: () => void;
}

const typeConfig = {
  need: {
    label: 'Need',
    description: 'Essential expenses',
    color: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
    icon: '🏠'
  },
  want: {
    label: 'Want',
    description: 'Discretionary spending',
    color: 'bg-purple-500/10 text-purple-500 border-purple-500/20',
    icon: '✨'
  },
  savings: {
    label: 'Savings',
    description: 'Money set aside',
    color: 'bg-green-500/10 text-green-500 border-green-500/20',
    icon: '💰'
  }
};

export function CategoryTypeSelector({ category, onUpdate }: CategoryTypeSelectorProps) {
  const [value, setValue] = useState<string>(category.category_type || '');
  const updateCategoryType = useUpdateCategoryType();

  const handleChange = async (newValue: string) => {
    const categoryType = newValue === 'none' ? undefined : (newValue as 'need' | 'want' | 'savings');

    try {
      await updateCategoryType.mutateAsync({
        id: category.id,
        data: { category_type: categoryType }
      });
      setValue(newValue === 'none' ? '' : newValue);
      onUpdate?.();
    } catch (error) {
      console.error('Failed to update category type:', error);
    }
  };

  const currentType = value as keyof typeof typeConfig;
  const config = currentType ? typeConfig[currentType] : null;

  return (
    <Select
      value={value || 'none'}
      onValueChange={handleChange}
      disabled={updateCategoryType.isPending}
    >
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select type">
          {config ? (
            <div className="flex items-center gap-2">
              <span>{config.icon}</span>
              <span>{config.label}</span>
            </div>
          ) : (
            <span className="text-muted-foreground">Unclassified</span>
          )}
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="none">
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground">-</span>
            <span className="text-muted-foreground">Unclassified</span>
          </div>
        </SelectItem>

        {Object.entries(typeConfig).map(([key, config]) => (
          <SelectItem key={key} value={key}>
            <div className="flex items-center gap-2">
              <span>{config.icon}</span>
              <div>
                <div>{config.label}</div>
                <div className="text-xs text-muted-foreground">{config.description}</div>
              </div>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

interface CategoryTypeBadgeProps {
  type?: 'need' | 'want' | 'savings';
}

export function CategoryTypeBadge({ type }: CategoryTypeBadgeProps) {
  if (!type || !typeConfig[type]) {
    return (
      <Badge variant="outline" className="text-muted-foreground">
        Unclassified
      </Badge>
    );
  }

  const config = typeConfig[type];

  return (
    <Badge variant="outline" className={cn(config.color)}>
      <span className="mr-1">{config.icon}</span>
      {config.label}
    </Badge>
  );
}
